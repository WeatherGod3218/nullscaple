package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/WeatherGod3218/nullscaple/enemies"
	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/nullscaple/logging"
	"github.com/sirupsen/logrus"
)

func GetEnemies(c *gin.Context) {
	loadedEnemies := enemies.GetEnemyList()
	c.JSON(http.StatusOK, loadedEnemies)
}

func GuessEnemy(c *gin.Context) {
	enemyId := c.PostForm("enemy-id")
	gameplayMode := c.PostForm("gameplay-mode")

	if !enemies.CheckIfStringIsMode(gameplayMode) {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to find the gamemode marked with the %d", gameplayMode))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	enemyIndex, err := strconv.Atoi(enemyId)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to convert the id %s", enemyId))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	foundEnemy := enemies.GetEnemyFromId(enemyIndex)
	if foundEnemy == nil {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to find enemy with the id %d", enemyIndex))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	baseEnemy := enemies.GetEnemyOfTheDay(gameplayMode)
	guessResult := enemies.CompareEnemies(*foundEnemy, *baseEnemy)

	c.JSON(http.StatusOK, gin.H{
		"result": guessResult,
	})
}

func GetTodaysEnemy(c *gin.Context) {
	gameplayMode := c.PostForm("gameplay-mode")
	if !enemies.CheckIfStringIsMode(gameplayMode) {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to find the gamemode marked with the %d", gameplayMode))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	baseEnemy := enemies.GetEnemyOfTheDay(gameplayMode)
	c.JSON(http.StatusOK, gin.H{
		"enemy": baseEnemy,
	})
}
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "Okay",
	})
}
