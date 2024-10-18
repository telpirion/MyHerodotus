package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	projectID  string
	userEmail  string = "anonymous@example.com"
)

type ClientError struct {
	Code       string `json:"code" binding:"required"`
	Message    string `json:"message" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Credential string `json:"credential" binding:"required"`
}

func main() {

	projectID = os.Getenv("PROJECT_ID")
	if _loggerName, ok := os.LookupEnv("LOGGER_NAME"); ok {
		loggerName = _loggerName
	}
	LogInfo("Starting Herodotus...")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/js", "./js")
	r.StaticFile("/favicon.ico", "./favicon.ico")

	r.GET("/home", startConversation)
	r.POST("/home", respondToUser)
	r.GET("/", login)
	r.POST("/logClientError", clientError)
	r.GET("/error", errPage)
	log.Fatal(r.Run(":8080"))
}

func login(c *gin.Context) {
	LogInfo("Login request received")
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func startConversation(c *gin.Context) {
	params := c.Params
	log.Println(params)
	LogInfo("Start conversation request received")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Message": struct {
			Message string
		}{
			Message: "Hello! I hear that you want to go on a trip somewhere. Tell me about it.",
		},
	})
}

func respondToUser(c *gin.Context) {
	LogInfo("Respond to user request received")

	c.Request.ParseForm()
	userMsg := c.Request.Form["userMsg"][0]
	log.Println(userMsg)

	botResponse, err := textPredictGemini(userMsg, projectID)
	if err != nil {
		LogError(fmt.Sprintf("Bad response from Gemini  %v\n", err))
		botResponse = "Oops! I had troubles understanding that ..."
	}

	convo := &ConversationBit{
		UserQuery:   userMsg,
		BotResponse: botResponse,
		Created:     time.Now(),
	}

	// Use a separate thread to store the conversation
	go func() {
		err := saveConversation(*convo, userEmail, projectID)
		if err != nil {
			LogError(fmt.Sprintf("Couldn't save conversation: %v\n", err))
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

func errPage(c *gin.Context) {
	c.HTML(http.StatusOK, "error.html", gin.H{})
}

func clientError(c *gin.Context) {
	var cError ClientError
	if err := c.ShouldBindJSON(&cError); err != nil {
		LogError(fmt.Sprintf("clientError JSON: %v\n", err))
	}
	LogError(fmt.Sprintf("clientError: %s, %s, %s\n", cError.Code, cError.Message, cError.Email))
	c.JSON(http.StatusOK, gin.H{"error": "message logged"})
}
