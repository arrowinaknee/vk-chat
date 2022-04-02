package server

import "net/http"

type Server struct {
	Addr   string
	DBLink string

	db  db_t
	msg messages
}

func (s *Server) Start() {
	s.db = db_t(s.DBLink)

	s.routes()
	http.ListenAndServe(s.Addr, nil)
}
