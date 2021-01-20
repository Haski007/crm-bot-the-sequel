package model

import "time"

type Cash struct {
	Title     string    `json:"title" bson:"title"`
	Amount    Money     `json:"amount" bson:"amount"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
