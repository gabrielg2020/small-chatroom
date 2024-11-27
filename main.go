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

	// Define websocket route
	router.GET("/ws", func(ctx *gin.Context) {
		// Serve the websocket connection
		ServeWs(hub, ctx.Writer, ctx.Request)
	})

	// Sever entry point
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})


	// Run the server
	router.Run(":8080")
}
