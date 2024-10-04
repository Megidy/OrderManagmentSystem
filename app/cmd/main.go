package main

import (
	"log"

	"github.com/Megidy/OrderManagmentSystem/pkg/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	routes.InitRoutes(router)
	log.Println("started server on port 8080")
	router.Run(":8080")
}
