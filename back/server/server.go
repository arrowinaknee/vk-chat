package server

import (
	"net/http"

	"ru.arrowinaknee.vk-chat/db"
)

type Server struct {
	Addr   string
	DBLink string
	DBName string

	db  db.DB
	msg messages
}

func (s *Server) Start() {
	s.db = db.DB{
		Uri: s.DBLink,
		Use: s.DBName,
	}

	s.routes()
	http.ListenAndServe(s.Addr, nil)
}
