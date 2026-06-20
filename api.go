package main

import (
	"fmt"
	"net/http"

	"github.com/WeatherGod3218/nullscaple/internal/database"
	"github.com/WeatherGod3218/nullscaple/internal/enemies"
	t "github.com/WeatherGod3218/nullscaple/internal/nulltypes"
	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/nullscaple/internal/logging"
	"github.com/sirupsen/logrus"
)

func GetEnemies(c *gin.Context) {
	loadedEnemies := enemies.GetEnemyList()
	c.JSON(http.StatusOK, loadedEnemies)
}

func GuessEnemy(c *gin.Context) {
	var req t.EnemyRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	difficulty, err := t.ParseDifficulty(req.Mode)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "GuessEnemy"}).Warn(fmt.Sprintf("Failed to find the gamemode marked with the %s", req.Mode))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	foundEnemy := enemies.GetEnemyFromId(req.ID)
	if foundEnemy == nil {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "GuessEnemy"}).Warn(fmt.Sprintf("Failed to find enemy with the id %d", req.ID))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	val, exists := c.Get(database.USER_COOKIE)
	playerId, ok := val.(string)
	if !exists || !ok {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "GuessEnemy"}).Warn(fmt.Sprintf("Failed to find the player with the id %d", playerId))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	baseEnemy, err := enemies.GetEnemyOfTheDay(req.Mode)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GuessEnemy"}).Warn("Failed to get the enemy of the day!")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	correct, results := enemies.CompareEnemies(foundEnemy, baseEnemy)

	remaining, err := database.AddPlayerGuess(playerId, req.ID, difficulty)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GuessEnemy"}).Warn("Failed to add player guess to database!")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	var gamestatus string
	if correct {
		gamestatus = "Win"
		database.SetPlayerGameResult(playerId, "Win", difficulty)
	} else if remaining <= 0 {
		gamestatus = "Lose"
		database.SetPlayerGameResult(playerId, "Lose", difficulty)
	}
	c.JSON(http.StatusOK, gin.H{
		"game_status": gamestatus,
		"correct":     correct,
		"result":      results,
		"enemy":       foundEnemy,
	})
}

func GetTodaysEnemy(c *gin.Context) {
	var req t.EnemyRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !enemies.CheckIfStringIsMode(req.Mode) {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "GetTodaysEnemy"}).Warn(fmt.Sprintf("Failed to find the gamemode marked with the %s", req.Mode))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	baseEnemy, err := enemies.GetEnemyOfTheDay(req.Mode)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enemy": baseEnemy,
	})
}

func GetHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func GetGuessScreenPage(c *gin.Context) {
	mode := c.Param("mode")

	if !enemies.CheckIfStringIsMode(mode) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	enemyList := enemies.GetEnemyList()

	c.HTML(http.StatusOK, "guessing.tmpl", gin.H{
		"Enemies": enemyList,
		"Mode":    mode,
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "Okay",
	})
}
