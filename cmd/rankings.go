package cmd

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SubmitRanking(dbService *DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var submission RankingSubmission

		if err := c.ShouldBindJSON(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Validate voter exists
		voters := dbService.GetVoters()
		voterExists := false
		for _, voter := range voters {
			if voter.ID == submission.VoterID && voter.Active {
				voterExists = true
				break
			}
		}

		if !voterExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid or inactive voter ID",
			})
			return
		}

		// Validate teams
		teams := dbService.GetTeams()
		if len(submission.Teams) == 0 || len(submission.Teams) > len(teams) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid number of teams in ranking",
			})
			return
		}

		// Validate all team IDs exist
		teamMap := make(map[int]bool)
		for _, team := range teams {
			teamMap[team.ID] = true
		}

		for _, teamID := range submission.Teams {
			if !teamMap[teamID] {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid team ID: " + strconv.Itoa(teamID),
				})
				return
			}
		}

		// Submit ranking
		err := dbService.SubmitRanking(submission.VoterID, submission.Week, submission.Teams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to submit ranking",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":      "Ranking submitted successfully",
			"voter_id":     submission.VoterID,
			"week":         submission.Week,
			"teams_ranked": len(submission.Teams),
		})
	}
}

func GetRankings(dbService *DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get week parameter or default to current week
		weekStr := c.DefaultQuery("week", "1")
		week, err := strconv.Atoi(weekStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week parameter"})
			return
		}

		rankings := dbService.CalculateRankings(week)

		c.JSON(http.StatusOK, gin.H{
			"week":       week,
			"rankings":   rankings,
			"updated_at": time.Now(),
		})
	}
}

func GetRankingsByWeek(dbService *DatabaseService) gin.HandlerFunc {
	return func(c *gin.Context) {
		weekStr := c.Param("week")
		week, err := strconv.Atoi(weekStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week parameter"})
			return
		}

		rankings := dbService.CalculateRankings(week)

		c.JSON(http.StatusOK, gin.H{
			"week":       week,
			"rankings":   rankings,
			"updated_at": time.Now(),
		})
	}
}
