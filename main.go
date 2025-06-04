package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Team struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type Vote struct {
	ID       int       `json:"id"`
	VoterID  int       `json:"voter_id"`
	TeamID   int       `json:"team_id"`
	Rank     int       `json:"rank"`
	Week     int       `json:"week"`
	CreateAt time.Time `json:"created_at"`
}

type RankingSubmission struct {
	VoterID int   `json:"voter_id" binding:"required"`
	Week    int   `json:"week" binding:"required"`
	Teams   []int `json:"teams" binding:"required"`
}

type RankingResponse struct {
	TeamID     int    `json:"team_id"`
	TeamName   string `json:"team_name"`
	Rank       int    `json:"rank"`
	Points     int    `json:"points"`
	FirstVotes int    `json:"first_votes"`
}

// Mock data for development
var teams = []Team{
	{ID: 1, Name: "T1", Region: "LCK"},
	{ID: 2, Name: "Gen.G", Region: "LCK"},
	{ID: 3, Name: "JDG", Region: "LPL"},
	{ID: 4, Name: "BLG", Region: "LPL"},
	{ID: 5, Name: "G2", Region: "LEC"},
}

var votes []Vote

func main() {
	// Create Gin router
	r := gin.Default()

	// Middleware for CORS (if you plan to have a frontend)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/health", healthCheck)
		api.GET("/teams", getTeams)
		api.POST("/rankings", submitRanking)
		api.GET("/rankings", getRankings)
		api.GET("/rankings/week/:week", getRankingsByWeek)
	}

	// Start server
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Local server couldn't start")
	}
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "lol-rankings-api",
	})
}

func getTeams(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"teams": teams,
		"count": len(teams),
	})
}

func submitRanking(c *gin.Context) {
	var submission RankingSubmission

	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if len(submission.Teams) == 0 || len(submission.Teams) > len(teams) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid number of teams in ranking",
		})
		return
	}

	// Clear existing votes for this voter and week
	votes = removeVotes(votes, submission.VoterID, submission.Week)

	// Add new votes
	for rank, teamID := range submission.Teams {
		vote := Vote{
			ID:       len(votes) + 1,
			VoterID:  submission.VoterID,
			TeamID:   teamID,
			Rank:     rank + 1, // rank starts at 1, not 0
			Week:     submission.Week,
			CreateAt: time.Now(),
		}
		votes = append(votes, vote)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Ranking submitted successfully",
		"voter_id":     submission.VoterID,
		"week":         submission.Week,
		"teams_ranked": len(submission.Teams),
	})
}

// Get current rankings (aggregate)
func getRankings(c *gin.Context) {
	// Get week parameter or default to current week
	weekStr := c.DefaultQuery("week", "1")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week parameter"})
		return
	}

	rankings := calculateRankings(week)

	c.JSON(http.StatusOK, gin.H{
		"week":       week,
		"rankings":   rankings,
		"updated_at": time.Now(),
	})
}

// Get rankings for specific week
func getRankingsByWeek(c *gin.Context) {
	weekStr := c.Param("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week parameter"})
		return
	}

	rankings := calculateRankings(week)

	c.JSON(http.StatusOK, gin.H{
		"week":       week,
		"rankings":   rankings,
		"updated_at": time.Now(),
	})
}

func calculateRankings(week int) []RankingResponse {
	teamPoints := make(map[int]int)
	teamFirstVotes := make(map[int]int)

	for _, vote := range votes {
		if vote.Week == week {
			points := max(0, 26-vote.Rank)
			teamPoints[vote.TeamID] += points

			if vote.Rank == 1 {
				teamFirstVotes[vote.TeamID]++
			}
		}
	}

	// Convert to response format and sort
	var rankings []RankingResponse
	for _, team := range teams {
		if points, exists := teamPoints[team.ID]; exists {
			rankings = append(rankings, RankingResponse{
				TeamID:     team.ID,
				TeamName:   team.Name,
				Points:     points,
				FirstVotes: teamFirstVotes[team.ID],
			})
		}
	}

	// Sort by points (descending)
	for i := 0; i < len(rankings)-1; i++ {
		for j := i + 1; j < len(rankings); j++ {
			if rankings[j].Points > rankings[i].Points {
				rankings[i], rankings[j] = rankings[j], rankings[i]
			}
		}
	}

	// Assign ranks
	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	return rankings
}

// Helper function to remove existing votes
func removeVotes(allVotes []Vote, voterID, week int) []Vote {
	var filtered []Vote
	for _, vote := range allVotes {
		if !(vote.VoterID == voterID && vote.Week == week) {
			filtered = append(filtered, vote)
		}
	}
	return filtered
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
