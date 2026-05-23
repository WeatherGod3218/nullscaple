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
	var req enemies.EnemyRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !enemies.CheckIfStringIsMode(req.Mode) {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to find the gamemode marked with the %s", req.Mode))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	enemyIndex, err := strconv.Atoi(req.ID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to convert the id %s", req.ID))
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

	baseEnemy := enemies.GetEnemyOfTheDay(req.Mode)
	guessResult := enemies.CompareEnemies(foundEnemy, baseEnemy)

	c.JSON(http.StatusOK, gin.H{
		"result": guessResult,
		"enemy":  foundEnemy,
	})
}

func GetTodaysEnemy(c *gin.Context) {
	var req enemies.EnemyRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !enemies.CheckIfStringIsMode(req.Mode) {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "Guessenemy"}).Warn(fmt.Sprintf("Failed to find the gamemode marked with the %s", req.Mode))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	baseEnemy := enemies.GetEnemyOfTheDay(req.Mode)
	c.JSON(http.StatusOK, gin.H{
		"enemy": baseEnemy,
	})
}

func GetHomePage(c *gin.Context) {
	enemyList := enemies.GetEnemyList()

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Enemies": enemyList,
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "Okay",
	})
}
