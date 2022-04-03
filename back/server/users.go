package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ru.arrowinaknee.vk-chat/api"
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
	id := r.Id()

	conn, err := s.db.Connect()
	if err != nil {
		res.Error(http.StatusBadGateway, "database not available")
		return
	}
	c := conn.Database("arrowchat").Collection("users")

	if r.Id() == "" {
		var arr []*api_user

		var filter any
		if q := r.Query(); q != "" {
			filter = bson.D{{Key: "$or", Value: bson.A{
				bson.D{{Key: "_id", Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: q, Options: "i"}},
				}}},
				bson.D{{Key: "nickname", Value: bson.D{
					{Key: "$regex", Value: primitive.Regex{Pattern: q, Options: "i"}},
				}}},
			}}}
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
			var u db_user
			err = cur.Decode(&u)

			if err != nil {
				log.Printf("Error fetching user data: %s", err)
				res.Error(http.StatusInternalServerError, "could not fetch user data")
				return
			}

			arr = append(arr, &api_user{Login: u.Id, Nickname: u.Nickname})
		}

		res.Write(arr)

	} else {
		var u db_user
		err = c.FindOne(context.TODO(), bson.D{bson.E{Key: "_id", Value: id}}).Decode(&u)

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

func (s *Server) HandleUsersPost(res *api.JsonResponse, u *api_user) {
	conn, err := s.db.Connect()
	if err != nil {
		res.Error(http.StatusBadGateway, "database not available")
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
		res.Error(http.StatusBadRequest, "user already exists")
		return
	} else if err != nil {
		log.Printf("Cannot insert user into db: %s", err)
		res.Error(http.StatusInternalServerError, "could not add user")
		return
	}
}
