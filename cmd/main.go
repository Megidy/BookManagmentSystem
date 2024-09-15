package main

import (
	"github.com/Megidy/BookManagmentSystem/pkj/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.InitRoutes(router)
	router.Run(":8080")
}
