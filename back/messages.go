package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var connections []*websocket.Conn
var messages []string

func handleConn(ws *websocket.Conn) {
	connections = append(connections, ws)
	log.Printf("New connection, total: %d", len(connections))
	for {
		msg := make([]byte, 512)
		n, err := ws.Read(msg)
		if err != nil {
			if err == io.EOF {
				for i, o := range connections {
					if o == ws {
						connections = append(connections[:i], connections[i+1:]...)
					}
				}
				log.Printf("Connection closed, total: %d", len(connections))
				return
			}
			log.Fatal(err)
		}
		messages = append(messages, string(msg[:n]))
		for _, ows := range connections {
			fmt.Fprintf(ows, "%s\n", msg[:n])
		}
	}
}

func handleMessageHistory(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Printf("Error serving /messages/history: %s\n", err)
		return
	}
	w.Write(bytes)
}
