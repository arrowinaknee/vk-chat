package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type messages struct {
	connections []*websocket.Conn
	messages    []string
}

func (s *Server) handleConn(ws *websocket.Conn) {
	s.msg.connections = append(s.msg.connections, ws)
	log.Printf("New connection, total: %d", len(s.msg.connections))
	for {
		msg := make([]byte, 512)
		n, err := ws.Read(msg)
		if err != nil {
			if err == io.EOF {
				for i, o := range s.msg.connections {
					if o == ws {
						s.msg.connections = append(s.msg.connections[:i], s.msg.connections[i+1:]...)
					}
				}
				log.Printf("Connection closed, total: %d", len(s.msg.connections))
				return
			}
			log.Fatal(err)
		}
		s.msg.messages = append(s.msg.messages, string(msg[:n]))
		for _, ows := range s.msg.connections {
			fmt.Fprintf(ows, "%s\n", msg[:n])
		}
	}
}

func (s *Server) handleMessageHistory(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(s.msg.messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Printf("Error serving /messages/history: %s\n", err)
		return
	}
	w.Write(bytes)
}
