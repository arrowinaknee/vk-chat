package db

type User struct {
	Id       string `bson:"_id"`
	Nickname string `bson:"nickname"`
	RegDate  int64  `bson:"reg_date"`
}
