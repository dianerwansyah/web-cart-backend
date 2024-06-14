package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"UserID" json:"UserID"`
	Amount      int                `bson:"Amount" json:"Amount"`
	Created     time.Time          `bson:"Created" json:"Created"`
	LastUpdated time.Time          `bson:"LastUpdated" json:"LastUpdated"`
}

func (Coupon) TableName() string {
	return "coupons"
}
