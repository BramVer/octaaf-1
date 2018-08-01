package main

import (
	"octaaf/models"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gopkg.in/telegram-bot-api.v4"
)

type Web struct{}

func ginit() {

	var web Web

	if OctaafEnv == "production" {
		gin.SetMode("release")
	}

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.GET("/health", web.health)
	}

	kali := api.Group("/kali")
	{
		kali.GET("/quote", web.quote)
	}
	router.Run()
}

func (Web) health(c *gin.Context) {
	// HTTP status code, either 200 or 500
	status := 200
	// User facing status message, contains all services
	var statusMessage = gin.H{
		"redis":    "Ok",
		"postgres": "Ok",
		"telegram": "Ok",
	}

	// Redis
	redisErr := Redis.Ping().Err()
	if redisErr != nil {
		statusMessage["redis"] = redisErr
		status = 500
	}

	// Postgres
	postgresErr := DB.RawQuery("SELECT COUNT(pid) FROM pg_stat_activity;").Exec()
	if postgresErr != nil {
		statusMessage["postgres"] = postgresErr
		status = 500
	}

	// Telegram
	_, telegramErr := Octaaf.GetMe()
	if telegramErr != nil {
		statusMessage["telegram"] = telegramErr
		status = 500
	}

	c.JSON(status, statusMessage)
}

func (Web) quote(c *gin.Context) {
	quote := models.Quote{}
	err := DB.Where("chat_id = ?", KaliID).Order("random()").Limit(1).First(&quote)

	if err != nil {
		c.JSON(500, err)
	}

	config := tgbotapi.ChatConfigWithUser{
		ChatID:             KaliID,
		SuperGroupUsername: "",
		UserID:             quote.UserID}

	user, err := Octaaf.GetChatMember(config)

	quoteMap := structs.Map(quote)

	if err != nil {
		quoteMap["from"] = err
	} else {
		quoteMap["from"] = user
	}

	c.JSON(200, quoteMap)
}
