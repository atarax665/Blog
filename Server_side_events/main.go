package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// SSE endpoint that the clients will be listening to
	router.GET("/sse", func(c *gin.Context) {
		// Set the response header to indicate SSE content type
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// Allow all origins to access the endpoint (Else you will get CORS error)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		// Create a channel to send events to the client
		println("Client connected")
		eventChan := make(chan string)
		clients[eventChan] = struct{}{} // Add the client to the clients map
		defer func() {
			delete(clients, eventChan) // Remove the client when they disconnect
			close(eventChan)
		}()

		// Listen for client close and remove the client from the list
		notify := c.Writer.CloseNotify()
		go func() {
			<-notify
			fmt.Println("Client disconnected")
		}()

		// Continuously send data to the client
		for {
			data := <-eventChan
			println("Sending data to client", data)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
		}
	})

	// Handle POST request
	router.POST("/send-data", func(c *gin.Context) {
		data := c.PostForm("data")
		// print data to console
		println("Data received from client :", data)
		broadcast(data)
		c.JSON(http.StatusOK, gin.H{"message": "Data sent to clients"})
	})

	// Start the server
	err := router.Run(":3000")
	if err != nil {
		fmt.Println(err)
	}

}

// Clients is a list of channels to send events to connected clients
var clients = make(map[chan string]struct{})

// broadcast sends an event to all connected clients
func broadcast(data string) {
	for client := range clients {
		client <- data
	}
}
