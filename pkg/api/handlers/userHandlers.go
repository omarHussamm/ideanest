package handlers

import (
	"context"
	configs "demo/config"
	"demo/pkg/controllers"
	"demo/pkg/database/mongodb/repository"
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

	c.JSON(http.StatusOK, req)

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

	u, err := userRepo.GetUserByEmailAndPassword(ctx, req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, u)

}

func RefreshToken(c *gin.Context) {
}
