package client

import (
	"log"
	"net/http"

	"senec-monitor/types"

	"github.com/gorilla/websocket"
)

// this is just to test the websocket handler.
func ConnectToWs() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000/subscribe", http.Header{})
	if err != nil {
		log.Println("Error connection to websocket")
	}

	for {
		var data types.LocalApiDataWithCorrectTypesWithTimeStamp

		if err := conn.ReadJSON(&data); err != nil {
			log.Println("Error reading from the websocket", err)
		}
		log.Println(data)

	}
}
