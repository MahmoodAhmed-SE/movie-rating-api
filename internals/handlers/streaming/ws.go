package streaming

import (
	"log"
	"movie-rating-api-go/internals/services"
	"net/http"
	"time"

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

func WSConnHandler(w http.ResponseWriter, r *http.Request) {
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

	userId, err := services.ValidateWSToken(tknMsg.Token)

	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	GlobalBroker.Mu.Lock()
	subscriber := GlobalBroker.SubscribeNotificationReceiver("new_notification", userId)
	GlobalBroker.Mu.Unlock()

	defer func() {
		GlobalBroker.RemoveSubscriber("new_notification", *subscriber)
		conn.Close()
	}()

	// receiving notification messages (realtime)
loopLabel:
	for {
		select {
		case notification := <-subscriber.Ch:
			log.Printf("A notifiction has been sent: \n%v", notification)
			conn.WriteJSON(&notification)
		case <-time.After(time.Minute * 10):
			log.Println("Connection timed out:", userId)
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Timeout"))
			break loopLabel
		}
	}
}
