package server

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func (s *Server) routes() {
	ws_serv := websocket.Server{Handler: websocket.Handler(s.handleConn)}
	http.HandleFunc("/", s.handlePage)
	http.HandleFunc("/users", s.HandleUsers)
	http.HandleFunc("/messages/ws", ws_serv.ServeHTTP)
	http.HandleFunc("/messages/history", s.handleMessageHistory)
}
