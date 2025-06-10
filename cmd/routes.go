package cmd

import "github.com/gin-gonic/gin"

func SetupRoutes(api *gin.RouterGroup, dbService *DatabaseService) {
	// Health check
	api.GET("/health", HealthCheck)

	// Teams endpoints
	api.GET("/teams", GetTeams(dbService))

	// Voters endpoints
	api.GET("/voters", GetVoters(dbService))
	api.POST("/voters", CreateVoter(dbService))

	// Rankings endpoints
	api.POST("/rankings", SubmitRanking(dbService))
	api.GET("/rankings", GetRankings(dbService))
	api.GET("/rankings/week/:week", GetRankingsByWeek(dbService))
}
