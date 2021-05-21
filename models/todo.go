package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Todo struct
type Todo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Task        string             `json:"task" form:"task" binding:"required" bson:"task"`
	Student     string             `json:"student" form:"student" binding:"required" bson:"student"`
	IsCompleted bool               `json:"isCompleted" form:"isCompleted" bson:"isCompleted"`
	Version     int                `json:"version" bson:"__v"`
}
