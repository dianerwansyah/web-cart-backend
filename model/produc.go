package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string             `bson:"Name" json:"Name"`
	Description    string             `bson:"Description" json:"Description"`
	Price          float64            `bson:"Price" json:"Price"`
	ImageURL       string             `bson:"ImageURL" json:"ImageURL"`
	Stock          int                `bson:"Stock" json:"Stock"`
	CategoryID     []string           `bson:"CategoryID" json:"CategoryID"`
	Created        time.Time          `bson:"Created" json:"Created"`
	LastUpdate     time.Time          `bson:"LastUpdate" json:"LastUpdate"`
	LastUpdateByID primitive.ObjectID `bson:"LastUpdateByID" json:"LastUpdateByID"`
}

func (Product) TableName() string {
	return "products"
}
