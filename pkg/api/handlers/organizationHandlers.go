package handlers

import (
	"context"
	configs "demo/config"
	"demo/pkg/controllers"
	"demo/pkg/database/mongodb/models"
	"demo/pkg/database/mongodb/repository"
	"demo/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var organizationRepo = *repository.NewOrganizationRepository(configs.DB.Database("ideaNest"))

func CreateOrganization(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req controllers.NewOrganizationReq
	err := c.BindJSON(&req)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	name, email, err := utils.ExtractTokenClaims(c)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	organizationid, err := organizationRepo.CreateOrganization(ctx, req, models.OrganizationMember{Name: name, Email: email, AccessLevel: "admin"})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"organization_id": organizationid,
	})
}

func GetOrganization(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	organizationId := c.Param("organizationId")
	oid, err := strconv.ParseUint(organizationId, 10, 32)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	organization, err := organizationRepo.GetOrganizationByID(ctx, uint(oid))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, organization)
}

func GetAllOrganizations(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	organizations, err := organizationRepo.GetAllOrganizations(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, organizations)
}

func UpdateOrganization(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	organizationId := c.Param("organizationId")
	oid, err := strconv.ParseUint(organizationId, 10, 32)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var req controllers.NewOrganizationReq
	err = c.BindJSON(&req)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	_, email, err := utils.ExtractTokenClaims(c)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	res, err := organizationRepo.UpdateOrganizationByID(ctx, uint(oid), req, email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}
func DeleteOrganization(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	organizationId := c.Param("organizationId")
	oid, err := strconv.ParseUint(organizationId, 10, 32)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	_, email, err := utils.ExtractTokenClaims(c)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	err = organizationRepo.DeleteOrganizationByID(ctx, uint(oid), email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "organization deleted"})
}
func InviteUsertoOrganization(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	organizationId := c.Param("organizationId")
	oid, err := strconv.ParseUint(organizationId, 10, 32)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var req controllers.InviteReq
	err = c.BindJSON(&req)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	newMember_name, err := userRepo.GetUserNameByEmail(ctx, req.Email)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, email, err := utils.ExtractTokenClaims(c)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = organizationRepo.InviteUsertoOrganizationbyEmail(ctx, uint(oid), email, req.Email, newMember_name)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user added to the organization"})
}
