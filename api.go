package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/WeatherGod3218/nullscaple/redis"

	"github.com/WeatherGod3218/nullscaple/enemies"
	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/nullscaple/logging"
	"github.com/sirupsen/logrus"
)

func RedisRateLimiter(rate float64, capacity float64) gin.HandlerFunc {

	limiter := redis.NewTokenBucket(rate, capacity)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "RedisRateLimiter"}).Info(fmt.Sprintf("Recieved from IP:%s", ip))

		allowed, tokens, err := limiter.Allow(c, ip)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "RedisRateLimiter"}).Warn(fmt.Sprintf("Failure in the redis cache %v", err))
		} else if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests!",
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%.0f", capacity))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%v", tokens))

		c.Next()
	}
}

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
