package server

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
	"ru.arrowinaknee.vk-chat/api"
)

type messages struct {
	connections []*websocket.Conn
	history     []string
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
		s.msg.history = append(s.msg.history, string(msg[:n]))
		for _, ows := range s.msg.connections {
			fmt.Fprintf(ows, "%s\n", msg[:n])
		}
	}
}

func (s *Server) handleMessageHistory(res *api.JsonResponse, r *api.GetRequest) {
	res.Write(s.msg.history)
}
