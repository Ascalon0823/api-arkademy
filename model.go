package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	Password     string             `bson:"password"`
	PlayerRecord PlayerRecord       `bson:"playerRecord"`
}

type PlayerRecord struct {
	Characters []CharacterRecord
}

type CharacterRecord struct {
	CreationTime   primitive.DateTime
	LastPlayedTime primitive.DateTime
	Data           primitive.M
}
