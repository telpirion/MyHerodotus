package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	ai "github.com/telpirion/MyHerodotus/ai"
	db "github.com/telpirion/MyHerodotus/databases"
	"github.com/telpirion/MyHerodotus/generated"

	"github.com/gin-gonic/gin"
)

var (
	r              *gin.Engine
	projectID      string
	userEmail      string = "anonymous@example.com"
	encryptedEmail string
	userEmailParam string = "user"
	contextTokens  int32
)

type ClientError struct {
	Code       string `json:"code" binding:"required"`
	Message    string `json:"message" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Credential string `json:"credential" binding:"required"`
}

type UserRating struct {
	BotResponse string `json:"response" binding:"required"`
	UserRating  string `json:"rating"   binding:"required"`
	DocumentID  string `json:"document" binding:"required"`
}

type UserMessage struct {
	Message string `json:"message" binding:"required"`
	Model   string `json:"model" binding:"required"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	defer func() {
		if r := recover(); r != nil {
			LogError(fmt.Sprintf("Error: %v", r))
		}
	}()

	projectID = os.Getenv("PROJECT_ID")
	if _loggerName, ok := os.LookupEnv("LOGGER_NAME"); ok {
		loggerName = _loggerName
	}
	writeTimeSeriesValue(projectID, "Herodotus warming up")
	defer func() {
		writeTimeSeriesValue(projectID, "Herodotus shutting down")
	}()

	LogInfo("Starting Herodotus...")

	r = gin.Default()
	r.LoadHTMLGlob("../site/html/*")
	r.Static("/js", "../site/js")
	r.Static("/css", "../site/css")
	r.StaticFile("/favicon.ico", "./favicon.ico")

	r.GET("/home", startConversation)
	r.POST("/home", respondToUser)
	r.GET("/", login)
	r.POST("/logClientError", clientError)
	r.GET("/error", errPage)
	r.POST("/rateResponse", rateResponse)
	r.POST("/predict", predict)

	log.Fatal(r.Run(":8080"))
}

func login(c *gin.Context) {
	LogInfo("Login request received")
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func startConversation(c *gin.Context) {
	writeTimeSeriesValue(projectID, "Start of conversation")
	// extractParams will redirect if user isn't logged in.
	userEmail = extractParams(c)
	encryptedEmail = userEmail

	if os.Getenv("CONFIGURATION_NAME") != "HerodotusDev" {
		encryptedEmail = transformEmail(userEmail)
	}

	LogInfo("Start conversation request received")

	// create a new conversation context
	convoHistory, err := db.GetConversation(encryptedEmail, projectID)
	if err != nil {
		LogError(fmt.Sprintf("couldn't get conversation history: %v\n", err))
	}

	// VertexAI + Gemini caching has a hard lower minimum; warn if the
	// minimum isn't reached
	convoContext, err := ai.StoreConversationContext(convoHistory, projectID)
	var minConvoNum *ai.MinCacheNotReachedError
	if errors.As(err, &minConvoNum) {
		LogWarning(err.Error())
	} else if err != nil {
		LogError(fmt.Sprintf("couldn't store conversation context: %v\n", err))
	}

	contextTokens, err = ai.GetTokenCount(convoContext, projectID)
	if err != nil {
		LogWarning(fmt.Sprintf("couldn't get context token count: %v\n", err))
	}

	// Populate the conversation context variable for grounding both Gemma and
	// Gemini (< 33000 tokens) caching.
	err = ai.SetConversationContext(convoHistory)
	if err != nil {
		LogError(fmt.Sprintf("couldn't set conversation context: %v\n", err))
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Message": "Hello! I hear that you want to go on a trip somewhere. Tell me about it.",
	})
}

func respondToUser(c *gin.Context) {
	defer writeTimeSeriesValue(projectID, "End of conversation")
	// extractParams will redirect if user isn't logged in.
	userEmail = extractParams(c)

	LogInfo("Respond to user request received")

	// Parse data
	var userMsg UserMessage
	var botResponse string
	var promptTemplateName string
	err := c.BindJSON(&userMsg)
	if err != nil {
		responseMsg := fmt.Sprintf("could not parse client message: %v\n", err)
		LogError(responseMsg)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": responseMsg,
		})
		return
	}

	botResponse, promptTemplateName, err = ai.Predict(userMsg.Message, userMsg.Model, projectID)
	if err != nil {
		LogError(fmt.Sprintf("bad response from %s: %v\n", userMsg.Model, err))
		c.JSON(http.StatusOK, gin.H{
			"Message": "Oops! I had troubles understanding that ...",
		})
		return
	}

	// Store data in Firestore
	documentID, err := updateDatabase(projectID, userMsg.Message, userMsg.Model, promptTemplateName, botResponse)
	if err != nil {
		LogError(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": struct {
			Message    string
			Email      string
			DocumentID string
		}{
			Message:    botResponse,
			Email:      userEmail,
			DocumentID: documentID,
		},
	})
}

func updateDatabase(projectID, userMessage, modelName, promptTemplateName, botResponse string) (string, error) {
	// Remove PII from message
	cleanMsg, err := deidentify(projectID, userMessage)
	if err != nil {
		return "", fmt.Errorf("couldn't deidentify user message: %v", err)
	}

	// Remove any sensitive data from botResponse.
	botResponse, err = deidentify(projectID, botResponse)
	if err != nil {
		return "", fmt.Errorf("couldn't deidentify bot response: %v", err)
	}

	// Get the number of tokens in the user message and response
	botTokens, err := ai.GetTokenCount(botResponse, projectID)
	if err != nil {
		return "", fmt.Errorf("can't get bot token count: %v", err)
	}

	userTokens, err := ai.GetTokenCount(cleanMsg, projectID)
	if err != nil {
		return "", fmt.Errorf("can't get bot token count: %v", err)
	}

	convo := &generated.ConversationBit{
		UserQuery:   cleanMsg,
		BotResponse: botResponse,
		Created:     time.Now().Unix(),
		Model:       modelName,
		Prompt:      promptTemplateName,
		TokenCount:  botTokens + userTokens,
	}

	// Store the conversation in Firestore and update the cachedContext
	// This is dual-entry accounting so that we don't have to query Firestore
	// every time to update the cached context
	documentID, err := db.SaveConversation(*convo, encryptedEmail, projectID)
	if err != nil {
		return "", fmt.Errorf("couldn't save conversation: %v", err)
	}
	return documentID, nil
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

func extractParams(c *gin.Context) string {
	// Verify that the user is signed in before answering
	userEmail = c.Query(userEmailParam)
	if userEmail == "" {
		LogWarning("User attempted to navigate directly to /home; redirected to sign-in")
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	}
	return userEmail
}

func rateResponse(c *gin.Context) {
	LogInfo("User rating received")
	var userRating UserRating
	err := c.BindJSON(&userRating)
	if err != nil {
		LogError(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "JSON incorrect",
		})
		return
	}

	err = db.UpdateConversation(userRating.DocumentID, encryptedEmail, userRating.UserRating, projectID)
	if err != nil {
		LogError(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Couldn't record rating",
		})
		return
	}
	LogInfo("User rating successfully recorded")

	c.JSON(http.StatusOK, gin.H{"success": "rating logged"})
}

func predict(c *gin.Context) {
	var userMsg UserMessage
	err := c.BindJSON(&userMsg)
	if err != nil {
		responseMsg := fmt.Sprintf("predict: could not parse client message: %v\n", err)
		LogError(responseMsg)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": responseMsg,
		})
		return
	}

	botResponse, _, err := ai.Predict(userMsg.Message, userMsg.Model, projectID)
	if err != nil {
		responseMsg := fmt.Sprintf("predict: could not get prediction: %v\n", err)
		LogError(responseMsg)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": responseMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": struct {
			Message string
		}{
			Message: botResponse,
		},
	})
}
