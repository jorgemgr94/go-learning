package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var users = []string{"John", "Jane", "Jim", "Jill"}

func GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("getting all users")

		c.JSON(http.StatusOK, users)
	}
}

func GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		slog.Info("getting a user", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		slog.Info("getting user by id")

		c.JSON(http.StatusOK, users[intId])
	}
}

func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("creating a user")
		users = append(users, "John Doe")
		c.JSON(http.StatusOK, "user created")
	}
}
