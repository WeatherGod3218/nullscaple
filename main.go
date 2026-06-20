package main

import (
	"net/http"

	"github.com/WeatherGod3218/nullscaple/internal/database"
	"github.com/WeatherGod3218/nullscaple/internal/enemies"
	"github.com/WeatherGod3218/nullscaple/internal/logging"
	"github.com/WeatherGod3218/nullscaple/internal/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	enemies.InitEnemies()
	err := redis.InitRedis()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "InitRedis"}).Fatal("Failed to init redis!")
	}

	err = database.InitDatabase()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "InitDatabase"}).Fatal("Failed to init database!")
	}

	router := gin.Default()

	router.GET("/health", HealthCheck)

	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")

	users := router.Group("")
	users.Use(database.DatabaseCookie())
	users.GET("/", redis.RedisRateLimiter(2, 50), GetHomePage)

	users.GET("/guess-screen/:mode", redis.RedisRateLimiter(2, 50), GetGuessScreenPage)
	users.GET("/get-enemies", redis.RedisRateLimiter(2, 50), GetEnemies)
	users.GET("/get-todays-enemy", redis.RedisRateLimiter(2, 50), GetTodaysEnemy)
	users.GET("/player-guesses", redis.RedisRateLimiter(2, 50), GetPlayerGuesses)
	users.POST("/guess-enemy", redis.RedisRateLimiter(2, 50), GuessEnemy)

	router.Run(":8080")
}
