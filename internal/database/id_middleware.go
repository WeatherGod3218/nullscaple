package database

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/nullscaple/internal/idgen"
	"github.com/WeatherGod3218/nullscaple/internal/logging"
	"github.com/sirupsen/logrus"
)

const USER_COOKIE string = "UserId"

func HandleNewUser(ctx *gin.Context) (string, error) {
	newId, err := idgen.GenerateNewId()

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to generate a new UUID7")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return "", err
	}

	domain := os.Getenv("SERVER_URL")
	secure := os.Getenv("GIN_MODE") == "release"

	ctx.SetCookie(USER_COOKIE, newId, 2147483647, "/", domain, secure, true)

	err = CreatePlayerData(newId)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "id": newId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to Create Player Data")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return "", nil
	}

	return newId, nil
}

func HandleReturningUser(ctx *gin.Context, userId string) (string, error) {
	exists, err := CheckPlayerExists(userId)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "id": userId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to see if a player exists!")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return "", err
	}

	if !exists {
		err = CreatePlayerData(userId)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "id": userId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to Create Player Data")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return "", err
		}
	}

	return userId, nil
}

func DatabaseCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := ctx.Cookie(USER_COOKIE)

		if err != nil {
			userId, err = HandleNewUser(ctx)
			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "id": userId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to Handle New User!")
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			userId, err = HandleReturningUser(ctx, userId)
			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "id": userId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to Handle Returning User!")
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		refreshed, err := RefreshPlayer(userId)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "id": userId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to Check Refresh for Player")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if refreshed {
			//TODO
		}

		domain := os.Getenv("SERVER_URL")
		secure := os.Getenv("GIN_MODE") == "release"

		ctx.SetCookie(USER_COOKIE, userId, 2147483647, "/", domain, secure, true)

		ctx.Set(USER_COOKIE, userId)
		ctx.Next()
	}
}
