package server

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"ru.arrowinaknee.vk-chat/api"
	"ru.arrowinaknee.vk-chat/db"
)

// creates a token that identifies the user
// storing sessions would be a better option, but it can be implemented later
func (s *Server) makeToken(login string, pass string) (string, error) {
	conn, err := s.db.Connect()
	if err != nil {
		return "", err
	}
	c := conn.Collection("users")

	count, err := c.CountDocuments(context.TODO(), db.ById(login))
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", mongo.ErrNoDocuments
	}

	return login, nil
}

func (s *Server) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := r.Header.Get("Authorization")

		res := api.JsonResponse{W: w}
		conn, err := s.db.Connect()
		if err != nil {
			res.Error(http.StatusBadGateway, "database not available")
			return
		}
		c := conn.Collection("users")

		count, err := c.CountDocuments(context.TODO(), db.ById(uid))
		if err != nil {
			res.Error(http.StatusInternalServerError, "something went wrong")
			return
		}
		if count == 0 {
			res.Error(http.StatusUnauthorized, "wrong password or username")
			return
		}

		h(w, r)
	}
}

// does not check if user exists at all, should be used in combination with Auth middleware
func (s *Server) getUserId(r *http.Request) string {
	uid := r.Header.Get("Authorization")
	return uid
}

type login_data struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

func (s *Server) HandleAuthPost(res *api.JsonResponse, r *api.JsonRequest[login_data]) {
	token, err := s.makeToken(r.V.Login, r.V.Pass)
	if err == mongo.ErrNoDocuments {
		res.Error(http.StatusUnauthorized, "incorrect login or password")
		return
	} else if err != nil {
		log.Printf("error generating token: %s", err)
		res.Error(http.StatusInternalServerError, "something went wrong")
		return
	}
	res.Write(map[string]string{"token": token})
}
