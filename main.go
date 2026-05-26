package main

import (
	"net/http"

	"github.com/WeatherGod3218/nullscaple/enemies"
	"github.com/WeatherGod3218/nullscaple/redis"

	"github.com/WeatherGod3218/nullscaple/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	enemies.InitEnemies()
	err := redis.InitRedis()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "InitRedis"}).Fatal("Failed to init redis!")
	}
	router := gin.Default()

	router.GET("/health", HealthCheck)

	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")

	router.GET("/", RedisRateLimiter(2, 50), GetHomePage)

	router.GET("/guess-screen/:mode", RedisRateLimiter(2, 50), GetGuessScreenPage)
	router.GET("/get-enemies", RedisRateLimiter(2, 50), GetEnemies)
	router.POST("/guess-enemy", RedisRateLimiter(2, 50), GuessEnemy)
	router.GET("/get-enemy-today", RedisRateLimiter(2, 50), GetTodaysEnemy)

	router.Run(":8080")
}
