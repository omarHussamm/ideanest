package handlers

import (
	"context"
	configs "demo/config"
	"demo/pkg/controllers"
	"demo/pkg/database/mongodb/repository"
	"demo/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var userRepo = *repository.NewUserRepository(configs.DB.Database("ideaNest"))

func GetAllUsers(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	u, err := userRepo.GetAllUsers(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, u)

}

func SignUp(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req controllers.SignUpReq
	err := c.BindJSON(&req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = userRepo.CreateUser(ctx, req)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "welcome!"})

}

func SignIn(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req controllers.LoginReq
	err := c.BindJSON(&req)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	token, err := userRepo.GetUserByEmailAndPassword(ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "welcome back!",
		"access_token":  token["access_token"],
		"refresh_token": token["refresh_token"],
	})

}

func RefreshTokens(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tokenReq controllers.TokenReqBody
	c.BindJSON(&tokenReq)

	uid, err := utils.RefreshTokenValid(tokenReq.Refresh_token)

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userRepo.GetUserByID(ctx, uid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tokens refreshed", "access_token": token["access_token"], "refresh_token": token["refresh_token"]})

}
