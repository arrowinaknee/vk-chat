package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ru.arrowinaknee.vk-chat/api"
	"ru.arrowinaknee.vk-chat/db"
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

func (s *Server) HandleUsersGet(res *api.JsonResponse, r *api.GetRequest) {
	conn, err := s.db.Connect()
	if err != nil {
		res.Error(http.StatusBadGateway, "database not available")
		return
	}
	c := conn.Collection("users")

	if r.Id() == "" {
		var arr []*api_user

		var filter any
		if q := r.Query(); q != "" {
			filter = db.Or([]any{
				db.Regex("nickname", q, "i"),
				db.Regex("_id", q, "i"),
			})
		} else {
			filter = bson.D{}
		}

		cur, err := c.Find(context.TODO(), filter, options.Find())

		if err != nil {
			log.Printf("Error fetching user data: %s", err)
			res.Error(http.StatusInternalServerError, "could not fetch user data")
			return
		}

		for cur.Next(context.TODO()) {
			u, err := db.Decode[db.User](cur)

			if err != nil {
				log.Printf("Error fetching user data: %s", err)
				res.Error(http.StatusInternalServerError, "could not fetch user data")
				return
			}

			arr = append(arr, &api_user{Login: u.Id, Nickname: u.Nickname})
		}

		res.Write(arr)

	} else {
		u, err := db.Decode[db.User](c.FindOne(context.TODO(), db.ById(r.Id())))

		if err == mongo.ErrNoDocuments {
			res.NotFound()
			return
		} else if err != nil {
			log.Printf("Error fetching user data: %s", err)
			res.Error(http.StatusInternalServerError, "could not fetch user data")
			return
		}

		res.Write(&api_user{Login: u.Id, Nickname: u.Nickname})
	}
}

func (s *Server) HandleUsersPost(res *api.JsonResponse, r *api.JsonRequest[api_user]) {
	conn, err := s.db.Connect()
	if err != nil {
		res.Error(http.StatusBadGateway, "database not available")
		return
	}
	c := conn.Collection("users")
	t := time.Now().Unix()
	_, err = c.InsertOne(context.TODO(), db_user{
		Id:       r.V.Login,
		Nickname: r.V.Nickname,
		RegDate:  t,
	})

	if mongo.IsDuplicateKeyError(err) {
		res.Error(http.StatusBadRequest, "user already exists")
		return
	} else if err != nil {
		log.Printf("Cannot insert user into db: %s", err)
		res.Error(http.StatusInternalServerError, "could not add user")
		return
	}
}
