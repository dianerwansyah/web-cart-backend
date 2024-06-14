package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Role     string             `bson:"role" json:"role"`
	Created  time.Time          `bson:"created" json:"created"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}
