package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Create a new hub
	hub := NewHub()
	// Run hub in a goroutine
	go hub.Run()

	// Sever entry point
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})


	// Run the server
	router.Run(":8080")
}
