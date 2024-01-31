package repository

import (
	"context"
	"demo/pkg/controllers"
	"demo/pkg/database/mongodb/models"
	"demo/pkg/utils"
	"fmt"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	err := r.GetUserByEmail(ctx, req.Email)

	if err == nil {
		return fmt.Errorf("email used before")
	}

	newUser := models.User{
		ID:       uint(rand.Uint32()),
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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) error {
	var user models.User
	err := r.collection.FindOne(ctx, gin.H{"email": email}).Decode(&user)
	if err != nil {
		return fmt.Errorf("error fetching user by email")
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, gin.H{"id": id}).Decode(&user)
	if err != nil {
		log.Printf("Error fetching user by id: %v\n", err)
		return nil, fmt.Errorf("error fetching user by id %v", err)
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmailAndPassword(ctx context.Context, req controllers.LoginReq) (map[string]string, error) {
	var user models.User
	err := r.collection.FindOne(ctx, req).Decode(&user)
	if err != nil {
		log.Printf("Error fetching user by email and password: %v\n", err)
		return nil, fmt.Errorf("error fetching user by email and password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Name)

	if err != nil {
		return nil, err
	}

	return token, nil
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

func (r *UserRepository) GetUserNameByEmail(ctx context.Context, email string) (string, error) {
	var user models.User
	err := r.collection.FindOne(ctx, gin.H{"email": email}).Decode(&user)
	if err != nil {
		log.Printf("Error fetching user by id: %v\n", err)
		return "", fmt.Errorf("error fetching user by id %v", err)
	}
	return user.Name, nil
}
