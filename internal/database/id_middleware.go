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
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			domain := os.Getenv("SERVER_URL")
			secure := os.Getenv("GIN_MODE") == "release"

			ctx.SetCookie(USER_COOKIE, newId, 2147483647, "/", domain, secure, true)
			userId = newId

			CreatePlayerData(newId)
		} else {
			exists, err := CheckPlayerExists(userId)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			if !exists {
				CreatePlayerData(userId)
			}
		}
		logging.Logger.WithFields(logrus.Fields{"id": userId, "module": "database", "method": "DatabaseCookie"}).Info("Got User Id:")

		ctx.Set(USER_COOKIE, userId)
		ctx.Next()
	}
}
