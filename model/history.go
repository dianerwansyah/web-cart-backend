package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type History struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	IDTrx      primitive.ObjectID `bson:"IDTrx" json:"IDTrx"`
	ProductID  primitive.ObjectID `bson:"ProductID" json:"ProductID"`
	UserID     primitive.ObjectID `bson:"UserID" json:"UserID"`
	Quantity   int                `bson:"Quantity" json:"Quantity"`
	IsCheckout bool               `bson:"IsCheckout" json:"IsCheckout"`
	IsConfirm  bool               `bson:"IsConfirm" json:"IsConfirm"`
	Created    time.Time          `bson:"Created" json:"Created"`
}

func (History) TableName() string {
	return "historys"
}
