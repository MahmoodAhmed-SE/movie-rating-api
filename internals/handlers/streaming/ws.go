package streaming

import (
	"fmt"
	"log"
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

var connTypes = map[string]struct{}{
	"notifications": {},
}

type Notification struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

func WSConnHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while calling Upgrade method: %v", err)
		return
	}
	i := time.Now().Unix()
	ii := i

	for {
		i = time.Now().Unix()
		if i >= ii {
			ii = time.Now().Add(time.Second * 5).Unix()

			notication := Notification{
				Date:    fmt.Sprint(i),
				Content: "New message!",
			}

			err = conn.WriteJSON(&notication)

			if err != nil {
				log.Printf("Error while writing json message: %v", err)
				return
			}
		}
	}
}
