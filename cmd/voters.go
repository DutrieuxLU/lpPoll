// cmd/voters.go
package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVoters(dbService *DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		voters := dbService.GetVoters()
		c.JSON(http.StatusOK, gin.H{
			"voters": voters,
			"count":  len(voters),
		})
	}
}

func CreateVoter(dbService *DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateVoterRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Check if email already exists
		voters := dbService.GetVoters()
		for _, voter := range voters {
			if voter.Email == request.Email {
				c.JSON(http.StatusConflict, gin.H{
					"error": "Voter with this email already exists",
				})
				return
			}
		}

		// Use default role if not provided
		role := request.Role
		if role == "" {
			role = "voter" // Match the database default
		}

		voter, err := dbService.CreateVoter(request.Name, request.Email, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create voter",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Voter created successfully",
			"voter":   voter,
		})
	}
}
