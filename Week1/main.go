package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", startConversation)
	r.POST("/", respondToUser)
	log.Fatal(r.Run(":8080"))
}

func startConversation(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Message": struct {
			Message string
		}{
			Message: "Hello! I hear that you want to go on a trip somewhere. Tell me about it.",
		},
	})
}

func respondToUser(c *gin.Context) {
	c.Request.ParseForm()
	userMsg := c.Request.Form["userMsg"][0]
	log.Println(userMsg)

	botResponse, err := textPredict(userMsg, "erschmid-test-291318", "us-central1", "text-bison@001")
	if err != nil {
		log.Println(err)
		botResponse = "Oops! I had troubles understanding that ..."
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Message": struct {
			Message string
		}{
			Message: botResponse,
		},
	})
}
