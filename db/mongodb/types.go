package mongodb

import (
	"github.com/ukane-philemon/megtask/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type dbUser struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt int64              `bson:"createdAt"`
}

type dbTask struct {
	ID          primitive.ObjectID `bson:"_id"`
	OwnerID     string             `bson:"ownerID"`
	db.TaskInfo `bson:"inline"`
}
