package main

import (
	"chat-rooms/client"

	"github.com/gin-gonic/gin"
)

func main() {

	h := client.NewHub()

	go h.Run()

	router := gin.New()
	router.LoadHTMLFiles("index.html")

	router.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		client.ServeWS(c.Writer, c.Request, roomId, h)
	})

	router.Run("0.0.0.0:8080")
}
