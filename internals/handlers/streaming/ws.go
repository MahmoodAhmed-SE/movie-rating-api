package streaming

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	constants "movie-rating-api-go/internals"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TO-DO: only allow acceptable client domains
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type FirstMessage struct {
	Token string `json:"jwt_token"`
}

func _validateWSToken(inputToken string) error {
	_, scanErr := fmt.Sscanf(inputToken, "Bearer %s", &inputToken)
	if scanErr != nil {
		log.Printf("Error scanning header authorization token: %v", scanErr)
		return scanErr
	}

	token, err := jwt.Parse(inputToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv(constants.EnvJWTSecretKey)), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return err
	}

	if !token.Valid {
		log.Println("invalid token")
		return errors.New("invalid token")
	}

	return nil
}

var stdinAdminOnce sync.Once

func WSConnHandler(w http.ResponseWriter, r *http.Request) {
	// simulating an admin sending new notifications to all conns through standard input
	stdinAdminOnce.Do(func() {
		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					break
				}
				for _, ch := range Connections {
					ch <- Notification{
						Date:    time.Now().Format("02/01/2006"),
						Content: string(line),
					}
				}
			}
		}()
	})

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while calling Upgrade method: %v", err)
		return
	}

	var tknMsg FirstMessage
	err = conn.ReadJSON(&tknMsg)

	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
			log.Println("Client closed connection before sending token.")
		} else {
			log.Printf("Error reading token message: %v", err)
		}
		conn.Close()
		return
	}

	if err = _validateWSToken(tknMsg.Token); err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	connId, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		conn.Close()
		return

	}

	Mu.Lock()
	Connections[connId] = make(chan Notification, 10)
	Mu.Unlock()

	defer func() {
		Mu.Lock()
		close(Connections[connId])
		delete(Connections, connId)
		Mu.Unlock()
	}()

	defer conn.Close()

	// receiving notification messages (realtime)
loopLabel:
	for {
		var ch chan Notification
		Mu.RLock()
		ch = Connections[connId]
		Mu.RUnlock()

		select {
		case notification := <-ch:
			log.Printf("A notifiction has been sent: \n%v", notification)
			conn.WriteJSON(&notification)
		case <-time.After(time.Minute * 10):
			log.Println("Connection timed out:", connId)
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Timeout"))
			break loopLabel
		}
	}
}
