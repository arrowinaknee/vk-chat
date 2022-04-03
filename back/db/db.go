package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Uri string
	Use string
}

func (db DB) Connect() (r *mongo.Database, err error) {
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(string(db.Uri)))
	if err != nil {
		log.Printf("Error connecting to db: %s", err)
		return
	}
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("Could not reach db: %s", err)
	}
	r = c.Database(db.Use)
	return
}
