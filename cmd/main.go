package main

import (
	configs "demo/config"
	"demo/pkg/api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	configs.ConnectDB()

	router := gin.Default()
	routes.UserRoute(router)

	log.Fatal(router.Run(":8080"))
}
