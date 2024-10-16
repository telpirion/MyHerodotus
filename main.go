package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var projectID string

const TestEmail string = "testemail@example.com"

func main() {

	projectID = os.Getenv("PROJECT_ID")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/home", startConversation)
	r.POST("/home", respondToUser)
	r.GET("/", login)
	log.Fatal(r.Run(":8080"))
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func startConversation(c *gin.Context) {
	params := c.Params
	log.Println(params)

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

	botResponse, err := textPredictGemma(userMsg, projectID)
	if err != nil {
		log.Println(err)
		botResponse = "Oops! I had troubles understanding that ..."
	}

	convo := &ConversationBit{
		UserQuery:   userMsg,
		BotResponse: botResponse,
		Created:     time.Now(),
	}

	// Use a separate thread to store the conversation
	go func() {
		err := saveConversation(*convo, TestEmail, projectID)
		if err != nil {
			log.Println(err)
		}
	}()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Message": struct {
			Message string
		}{
			Message: botResponse,
		},
	})
}
