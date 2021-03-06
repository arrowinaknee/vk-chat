package server

import (
	"net/http"

	"golang.org/x/net/websocket"
	"ru.arrowinaknee.vk-chat/api"
)

func (s *Server) routes() {
	ws_serv := websocket.Server{Handler: websocket.Handler(s.handleConn)}
	http.Handle("/messages/ws", ws_serv)

	http.Handle("/users", api.Endpoint{
		Get:  api.Url(s.HandleUsersGet),
		Post: api.Json(s.HandleUsersPost),
	})
	http.Handle("/friends", api.Endpoint{
		Preprocess: s.Auth,
		Get:        api.Url(s.HandleFriendsGet),
		Post:       api.Json(s.HandleFriendsPost),
		Delete:     api.Url(s.HandleFriendsDelete),
	})
	http.Handle("/auth", api.Endpoint{
		Post: api.Json(s.HandleAuthPost),
	})
	http.Handle("/messages/history", api.Endpoint{
		Get: api.Url(s.handleMessageHistory),
	})

	http.HandleFunc("/", s.handlePage)
}
