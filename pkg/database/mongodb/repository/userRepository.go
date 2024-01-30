package repository

import (
	"context"
	"demo/pkg/controllers"
	"demo/pkg/database/mongodb/models"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, req controllers.SignUpReq) error {
	_, err := r.GetUserByEmail(ctx, req.Email)

	if err == nil {
		log.Printf("This email %s is used before: %v\n", req.Email, err)
		return fmt.Errorf("error inserting user")
	}

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err = r.collection.InsertOne(ctx, newUser)
	if err != nil {
		log.Printf("Error inserting user: %v\n", err)
		return fmt.Errorf("error inserting user")
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Printf("Error fetching user by email: %v\n", err)
		return nil, fmt.Errorf("error fetching user by email")
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmailAndPassword(ctx context.Context, req controllers.LoginReq) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, req).Decode(&user)
	if err != nil {
		log.Printf("Error fetching user by email and password: %v\n", err)
		return nil, fmt.Errorf("error fetching user by email and password")
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Error fetching all users: %v\n", err)
		return nil, fmt.Errorf("error fetching all users")
	}

	for cur.Next(ctx) {
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			return nil, fmt.Errorf("error fetching all users")
		}
		users = append(users, &elem)
	}

	defer cur.Close(ctx)

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error fetching all users")
	}

	return users, nil
}
