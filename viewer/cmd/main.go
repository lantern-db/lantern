package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lantern-db/lantern/viewer/config"
	"log"
)

func main() {
	viewerConfig := config.Load()

	router := gin.Default()
	router.Static("/v1", "./viewer/static/v1")

	router.GET("/config", func(c *gin.Context) {
		c.JSON(200, viewerConfig)
	})
	if err := router.Run(":8080"); err != nil {
		log.Panicln(err)
	}
}
