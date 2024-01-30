package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty"`
	Name     string             `json:"name,omitempty" binding:"required"`
	Email    string             `json:"email,omitempty" binding:"required"`
	Password string             `json:"password,omitempty" binding:"required"`
}
