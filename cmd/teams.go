package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTeams(dbService *DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		teams := dbService.GetTeams()
		c.JSON(http.StatusOK, gin.H{
			"teams": teams,
			"count": len(teams),
		})
	}
}
