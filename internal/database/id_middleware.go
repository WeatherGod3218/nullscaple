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

func DatabaseCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := ctx.Cookie(USER_COOKIE)

		if err != nil {
			newId, err := idgen.GenerateNewId()

			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to generate a new UUID7")
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			domain := os.Getenv("SERVER_URL")
			secure := os.Getenv("GIN_MODE") == "release"

			ctx.SetCookie(USER_COOKIE, newId, 2147483647, "/", domain, secure, true)
			userId = newId

			err = CreatePlayerData(userId)
			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "id": userId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to Create Player Data")
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			exists, err := CheckPlayerExists(userId)
			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "id": userId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to see if a player exists!")
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if !exists {
				err = CreatePlayerData(userId)
				if err != nil {
					logging.Logger.WithFields(logrus.Fields{"error": err, "id": userId, "module": "database", "method": "DatabaseCookie"}).Warn("Failed to Create Player Data")
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return
				}
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

		ctx.Set(USER_COOKIE, userId)
		ctx.Next()
	}
}
