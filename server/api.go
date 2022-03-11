package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BotOpen(c *gin.Context) {
	NewBot(c)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func StartAPI() {
	config := GetConfig(".")

	r := gin.Default()

	v1 := r.Group("/b", gin.BasicAuth(gin.Accounts{
		config.RpcUser: config.RpcPassword,
	}))
	{
		v1.POST("/open", BotOpen)
	}

	r.Run(":" + config.RpcPort)
}
