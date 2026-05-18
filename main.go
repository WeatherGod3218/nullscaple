package main

import (
	"github.com/WeatherGod3218/nullscaple/enemies"
	"github.com/gin-gonic/gin"
)

func main() {

	enemies.InitEnemies()
	router := gin.Default()

	router.GET("/health", HealthCheck)

	// router.StaticFS("/static", http.Dir("static"))
	// router.LoadHTMLGlob("templates/*")

	router.GET("/get-enemies", GetEnemies)
	router.POST("/guess-enemy", GuessEnemy)
	router.GET("/get-enemy-today", GetTodaysEnemy)

	router.Run(":8080")
}
