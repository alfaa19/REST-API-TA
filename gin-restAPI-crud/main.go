package main

import (
	config "github.com/alfaa19/gin-restAPI-crud/config/mysql"
	"github.com/alfaa19/gin-restAPI-crud/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	r.GET("/stats", controller.GetAll)
	r.GET("/stats/:id", controller.GetById)
	r.POST("/stats", controller.Create)
	r.PUT("/stats/:id", controller.Update)
	r.DELETE("/stats/:id", controller.Delete)

	r.Run(":8081")
}
