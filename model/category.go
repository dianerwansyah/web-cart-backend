package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string             `bson:"name" json:"name"`
	Description    string             `bson:"description" json:"description"`
	Created        time.Time          `bson:"created" json:"created"`
	LastUpdate     time.Time          `bson:"last_update" json:"last_update"`
	LastUpdateByID primitive.ObjectID `bson:"last_update_by_id" json:"last_update_by_id"`
}

func (Category) TableName() string {
	return "categories"
}
