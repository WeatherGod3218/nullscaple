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

	enemyIndex, err := strconv.Atoi(enemyId)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to convert the id %s", enemyId))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to convert the id %s", enemyId),
		})
		return
	}

	foundEnemy := enemies.GetEnemyFromId(enemyIndex)
	if foundEnemy == nil {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to find enemy with the id %d", enemyIndex))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to find enemy with the id %d", enemyIndex),
		})
		return
	}

	baseEnemy := enemies.GetEnemyOfTheDay()
	guessResult := enemies.CompareEnemies(*foundEnemy, *baseEnemy)

	c.JSON(http.StatusOK, gin.H{
		"result": guessResult,
	})
}

func GetTodaysEnemy(c *gin.Context) {
	baseEnemy := enemies.GetEnemyOfTheDay()
	c.JSON(http.StatusOK, gin.H{
		"enemy": baseEnemy,
	})
}
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "Okay",
	})
}
