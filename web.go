package main

import "github.com/gin-gonic/gin"

var Web *gin.Engine

func ginit() {

	if OctaafEnv == "production" {
		gin.SetMode("release")
	}

	Web = gin.Default()
	Web.GET("/status", octaafStatus)
	Web.Run()
}

func octaafStatus(c *gin.Context) {
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
