package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db_t string

func (db db_t) Connect() (c *mongo.Client, err error) {
	c, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(string(db)))
	if err != nil {
		log.Printf("Error connecting to db: %s", err)
	}
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("Could not reach db: %s", err)
	}
	return
}
