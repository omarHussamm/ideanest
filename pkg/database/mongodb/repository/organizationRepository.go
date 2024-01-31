package repository

import (
	"context"
	"demo/pkg/controllers"
	"demo/pkg/database/mongodb/models"
	"fmt"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrganizationRepository struct {
	collection *mongo.Collection
}

func NewOrganizationRepository(db *mongo.Database) *OrganizationRepository {
	return &OrganizationRepository{
		collection: db.Collection("organizations"),
	}
}

func (r *OrganizationRepository) CreateOrganization(ctx context.Context, req controllers.NewOrganizationReq, user models.OrganizationMember) (uint, error) {

	err := r.GetOrganizationByName(ctx, req.Name)

	if err == nil {
		return 0, fmt.Errorf("name used before")
	}

	newOrganization := models.Organization{
		ID:                  uint(rand.Uint32()),
		Name:                req.Name,
		Description:         req.Description,
		OrganizationMembers: []models.OrganizationMember{user},
	}

	_, err = r.collection.InsertOne(ctx, newOrganization)
	if err != nil {
		log.Printf("Error inserting user: %v\n", err)
		return 0, fmt.Errorf("error inserting user")
	}
	return newOrganization.ID, nil
}

func (r *OrganizationRepository) GetOrganizationByName(ctx context.Context, name string) error {
	var organization models.Organization
	err := r.collection.FindOne(ctx, gin.H{"name": name}).Decode(&organization)
	if err != nil {
		return fmt.Errorf("error fetching organization by name")
	}
	return nil
}

func (r *OrganizationRepository) GetAllOrganizations(ctx context.Context) ([]*models.Organization, error) {
	var organizations []*models.Organization
	cur, err := r.collection.Find(ctx, gin.H{})
	if err != nil {
		log.Printf("Error fetching all organizations: %v\n", err)
		return nil, fmt.Errorf("error fetching all organizations")
	}

	for cur.Next(ctx) {
		var elem models.Organization
		err := cur.Decode(&elem)
		if err != nil {
			return nil, fmt.Errorf("error fetching all organizations")
		}
		organizations = append(organizations, &elem)
	}

	defer cur.Close(ctx)

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error fetching all organizations")
	}

	return organizations, nil
}

func (r *OrganizationRepository) GetOrganizationByID(ctx context.Context, oid uint) (*models.Organization, error) {
	var organization models.Organization
	err := r.collection.FindOne(ctx, gin.H{"id": oid}).Decode(&organization)
	if err != nil {
		log.Printf("Error organization by id: %v\n", err)
		return nil, fmt.Errorf("error organization by id")
	}

	return &organization, nil
}
func (r *OrganizationRepository) UpdateOrganizationByID(ctx context.Context, oid uint, req controllers.NewOrganizationReq, user_email string) (*controllers.UpdateOrganizationRes, error) {
	_, err := r.IsAuthorized(ctx, oid, user_email)
	if err != nil {
		return nil, err
	}

	var organization models.Organization
	err = r.collection.FindOneAndUpdate(ctx, gin.H{"id": oid}, gin.H{"$set": gin.H{"name": req.Name, "description": req.Description}}).Decode(&organization)
	if err != nil {
		log.Printf("Error updating organization: %v\n", err)
		return nil, fmt.Errorf("error updating organization")
	}
	return &controllers.UpdateOrganizationRes{ID: organization.ID, Name: req.Name, Description: req.Description}, nil
}
func (r *OrganizationRepository) DeleteOrganizationByID(ctx context.Context, oid uint, user_email string) error {
	_, err := r.IsAuthorized(ctx, oid, user_email)
	if err != nil {
		return err
	}

	var organization models.Organization
	err = r.collection.FindOneAndDelete(ctx, gin.H{"id": oid}).Decode(&organization)
	if err != nil {
		log.Printf("Error organization by id: %v\n", err)
		return fmt.Errorf("error organization by id")
	}

	return nil
}
func (r *OrganizationRepository) InviteUsertoOrganizationbyEmail(ctx context.Context, oid uint, user_email, newMember_email, newMember_name string) error {

	org, err := r.IsAuthorized(ctx, oid, user_email)
	if err != nil {
		return err
	}

	organizationMembers := org.OrganizationMembers
	for _, member := range organizationMembers {
		if member.Email == newMember_email {
			return fmt.Errorf("user is already a member of this organization")
		}
	}

	member := models.OrganizationMember{Name: newMember_name, Email: newMember_email, AccessLevel: "member"}
	org.OrganizationMembers = append(org.OrganizationMembers, member)

	var organization models.Organization
	err = r.collection.FindOneAndReplace(ctx, gin.H{"id": oid}, org).Decode(&organization)
	if err != nil {
		log.Printf("error adding user to organization: %v\n", err)
		return fmt.Errorf("error adding user to organization")
	}
	return nil
}

func (r *OrganizationRepository) IsAuthorized(ctx context.Context, oid uint, user_email string) (*models.Organization, error) {
	organization, err := r.GetOrganizationByID(ctx, oid)
	if err != nil {
		log.Printf("Error invalid organization id: %v\n", err)
		return nil, fmt.Errorf("error invalid organization id")
	}
	organizationMembers := organization.OrganizationMembers
	for _, member := range organizationMembers {
		if member.Email == user_email {
			if member.AccessLevel == "admin" {
				return organization, nil
			} else {
				return nil, fmt.Errorf("unauthorized access: not an admin")
			}
		}
	}
	return nil, fmt.Errorf("unauthorized access: not a member of this organization")
}
