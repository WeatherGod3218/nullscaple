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

	foundEnemy, err := enemies.GetEnemyFromId(req.ID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "GuessEnemy"}).Warn(fmt.Sprintf("Failed to find enemy with the id %s", req.ID))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	val, exists := c.Get(database.USER_COOKIE)
	playerId, ok := val.(string)
	if !exists || !ok {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "GuessEnemy"}).Warn(fmt.Sprintf("Failed to find the player with the id %s", playerId))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	baseEnemy, err := enemies.GetEnemyOfTheDay(difficulty)
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

	difficulty, err := t.ParseDifficulty(req.Mode)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetGuessScreenPage"}).Warn(fmt.Sprintf("Failed to find the gamemode marked with the %s", difficulty))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	baseEnemy, err := enemies.GetEnemyOfTheDay(difficulty)
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

func GetPlayerGuesses(c *gin.Context) {
	val, exists := c.Get(database.USER_COOKIE)
	playerId, ok := val.(string)
	if !exists || !ok {
		logging.Logger.WithFields(logrus.Fields{"module": "api", "method": "GetGuessScreenPage"}).Warn(fmt.Sprintf("Failed to find the player with the id %s", playerId))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	difficulty, err := t.ParseDifficulty(c.GetHeader("Difficulty"))
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetGuessScreenPage"}).Warn(fmt.Sprintf("Failed to find the gamemode marked with the %s", c.GetHeader("Difficulty")))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	baseEnemy, err := enemies.GetEnemyOfTheDay(difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	guessList, err := database.GetPlayerGuesses(playerId, difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to process the request!",
		})
		return
	}

	var enemyComparisons []t.EnemyGuess
	for _, id := range guessList {
		foundEnemy, err := enemies.GetEnemyFromId(id)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetGuessScreenPage"}).Warn(fmt.Sprintf("Failed to find guesses for player %s", playerId))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unable to process the request!",
			})
			return
		}

		_, comparison := enemies.CompareEnemies(foundEnemy, baseEnemy)

		enemyComparisons = append(enemyComparisons, t.EnemyGuess{Enemy: foundEnemy, ComparisonResult: comparison})
	}

	if enemyComparisons == nil {
		enemyComparisons = []t.EnemyGuess{}
	}

	c.JSON(http.StatusOK, enemyComparisons)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "Okay",
	})
}
