package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Regex(key string, pattern string, opts string) bson.D {
	return bson.D{
		{Key: key, Value: bson.D{
			{Key: "$regex", Value: primitive.Regex{Pattern: pattern, Options: opts}},
		}},
	}
}

func Or(args []any) bson.D {
	return bson.D{{
		Key:   "$or",
		Value: args,
	}}
}

// specify wich fields to return from query
func Include(fields ...string) bson.D {
	r := make(bson.D, 0, len(fields))
	for _, v := range fields {
		r = append(r, bson.E{
			Key:   v,
			Value: 1,
		})
	}
	return r
}

func ById(id any) bson.D {
	return bson.D{{
		Key:   "_id",
		Value: id,
	}}
}

type Decoder interface {
	Decode(any) error
}

func Decode[T any](d Decoder) (v *T, err error) {
	v = new(T)
	err = d.Decode(v)
	return
}
