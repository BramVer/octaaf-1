package web

import (
	"octaaf/models"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gobuffalo/pop"
	"gopkg.in/telegram-bot-api.v4"
)

type Connections struct {
	Octaaf      *tgbotapi.BotAPI
	Postgres    *pop.Connection
	Redis       *redis.Client
	KaliID      int64
	Environment string
}

var conn Connections

func New(c Connections) *gin.Engine {
	conn = c

	if conn.Environment == "production" {
		gin.SetMode("release")
	}

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.GET("/health", health)
	}

	kali := api.Group("/kali")
	{
		kali.GET("/quote", quote)
	}

	return router
}

func health(c *gin.Context) {
	// HTTP status code, either 200 or 500
	status := 200
	// User facing status message, contains all services
	var statusMessage = gin.H{
		"redis":    "Ok",
		"postgres": "Ok",
		"telegram": "Ok",
	}

	// Redis
	_, redisErr := conn.Redis.Ping().Result()
	if redisErr != nil {
		statusMessage["redis"] = redisErr
		status = 500
	}

	// Postgres
	postgresErr := conn.Postgres.RawQuery("SELECT COUNT(pid) FROM pg_stat_activity;").Exec()
	if postgresErr != nil {
		statusMessage["postgres"] = postgresErr
		status = 500
	}

	// Telegram
	_, telegramErr := conn.Octaaf.GetMe()
	if telegramErr != nil {
		statusMessage["telegram"] = telegramErr
		status = 500
	}

	c.JSON(status, statusMessage)
}

func quote(c *gin.Context) {
	quote := models.Quote{}
	err := conn.Postgres.Where("chat_id = ?", conn.KaliID).Order("random()").Limit(1).First(&quote)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	config := tgbotapi.ChatConfigWithUser{
		ChatID:             conn.KaliID,
		SuperGroupUsername: "",
		UserID:             quote.UserID}

	user, err := conn.Octaaf.GetChatMember(config)

	quoteMap := structs.Map(quote)

	if err != nil {
		quoteMap["from"] = err
	} else {
		quoteMap["from"] = user
	}

	c.JSON(200, quoteMap)
}
