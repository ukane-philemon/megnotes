package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type dbUser struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt int64              `bson:"createdAt"`
}
