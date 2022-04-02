package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type api_user struct {
	Login    string `json:"login"`
	Nickname string `json:"nickname"`
}

type db_user struct {
	Id       string `bson:"_id"`
	Nickname string `bson:"nickname"`
	RegDate  int64  `bson:"reg_date"`
}

func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.FormValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Request must have 'id' specified")
			return
		}
		conn, err := s.db.Connect()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Database not available")
			return
		}
		c := conn.Database("arrowchat").Collection("users")

		var u db_user
		err = c.FindOne(context.TODO(), bson.D{bson.E{Key: "_id", Value: id}}).Decode(&u)
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "User does not exist")
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch user")
			log.Printf("Error fetching user data: %s", err)
			return
		}

		err = json.NewEncoder(w).Encode(&api_user{Login: u.Id, Nickname: u.Nickname})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not fetch user")
			log.Printf("Error encoding user data: %s", err)
		}
	case http.MethodPost:
		var u api_user
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Print("Could not parse json body")
			return
		}
		conn, err := s.db.Connect()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Print("Database not available")
			return
		}
		c := conn.Database("arrowchat").Collection("users")
		t := time.Now().Unix()
		_, err = c.InsertOne(context.TODO(), db_user{
			Id:       u.Login,
			Nickname: u.Nickname,
			RegDate:  t,
		})

		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "User '%s' already exists", u.Login)
			return
		} else if err != nil {
			log.Printf("Cannot insert user into db: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Could not add user")
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Method %s not supported", r.Method)
	}

}
