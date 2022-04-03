package server

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ru.arrowinaknee.vk-chat/api"
	"ru.arrowinaknee.vk-chat/db"
)

func (s *Server) HandleFriendsGet(res *api.JsonResponse, r *api.UrlRequest) {
	uid := r.Base.Header.Get("Authorization")

	conn, err := s.db.Connect()
	if err != nil {
		res.Error(http.StatusBadGateway, "database not available")
		return
	}
	c := conn.Collection("users")

	u, err := db.Decode[db.User](c.FindOne(context.TODO(), db.ById(uid), options.FindOne().SetProjection(bson.D{{Key: "friends", Value: 1}})))

	if err == mongo.ErrNoDocuments {
		res.Error(http.StatusUnauthorized, "user does not exist")
		return
	} else if err != nil {
		log.Printf("Error fetching friends list: %s", err)
		res.Error(http.StatusInternalServerError, "error fetching friends")
		return
	}

	f := make([]*api_user, 0, len(u.Friends))
	cur, err := c.Find(context.TODO(), db.ById(bson.D{{Key: "$in", Value: u.Friends}}))
	if err != nil {
		log.Printf("Error fetching friends list: %s", err)
		res.Error(http.StatusInternalServerError, "error fetching friends")
		return
	}

	for cur.Next(context.TODO()) {
		u, err = db.Decode[db.User](cur)
		if err != nil {
			log.Printf("Error decoding friend: %s", err)
			res.Error(http.StatusInternalServerError, "error fetching friends")
			return
		}
		f = append(f, &api_user{Login: u.Id, Nickname: u.Nickname})
	}

	res.Write(f)
}

func (s *Server) HandleFriendsPost(res *api.JsonResponse, r *api.JsonRequest[api.IdRequest[string]]) {
	uid := r.Base.Header.Get("Authorization")

	conn, err := s.db.Connect()
	if err != nil {
		res.Error(http.StatusBadGateway, "database not available")
		return
	}
	c := conn.Collection("users")

	// check if there is a user with the specified id
	count, err := c.CountDocuments(context.TODO(), db.ById(r.V.Id))
	if err != nil {
		log.Printf("Error searching for user: %s", err)
		res.Error(http.StatusBadGateway, "database not available")
		return
	}
	if count == 0 {
		res.Error(http.StatusNotFound, "this user does not exist")
		return
	}

	// add the user to friends list (addToSet keeps values unique)
	upd, err := c.UpdateByID(context.TODO(), uid, bson.D{{Key: "$addToSet", Value: bson.D{{Key: "friends", Value: r.V.Id}}}})
	if err != nil {
		log.Printf("Error adding friend: %s", err)
		res.Error(http.StatusInternalServerError, "error adding friend")
		return
	}
	if upd.MatchedCount == 0 {
		res.Error(http.StatusUnauthorized, "user does not exist")
		return
	}
	// document was found but not updated, this means that the user was already in friends list
	if upd.ModifiedCount == 0 {
		res.Error(http.StatusConflict, "this user is already in friends")
		return
	}
}

func (s *Server) HandleFriendsDelete(res *api.JsonResponse, r *api.UrlRequest) {
	uid := r.Base.Header.Get("Authorization")

	conn, err := s.db.Connect()
	if err != nil {
		res.Error(http.StatusBadGateway, "database not available")
		return
	}
	c := conn.Collection("users")

	upd, err := c.UpdateByID(context.TODO(), uid, bson.D{{Key: "$pull", Value: bson.D{{Key: "friends", Value: r.Id()}}}})
	if err != nil {
		log.Printf("Error removing friend: %s", err)
		res.Error(http.StatusInternalServerError, "error removing friend")
		return
	}
	if upd.MatchedCount == 0 {
		res.Error(http.StatusUnauthorized, "user does not exist")
		return
	}
	if upd.ModifiedCount == 0 {
		res.Error(http.StatusNotFound, "this user is not a friend")
	}
}
