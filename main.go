// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"

	"lpPoll/cmd"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "lppoll-freedb.cz4oskg6qyk9.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "Fisher2019,PULL"
	dbname   = "postgres"
)

func main() {
	// Initialize database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB: ", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Initialize the database service
	dbService := cmd.NewDatabaseService(db)

	// Load initial data from database
	err = dbService.LoadFromDB()
	if err != nil {
		log.Fatal("UNABLE TO LOAD FROM REMOTE DB: ", err)
	}

	// Create Gin router
	r := gin.Default()

	// Setup CORS middleware
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

	// Setup API routes
	api := r.Group("/api/v1")
	cmd.SetupRoutes(api, dbService)

	// Start server
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Local server couldn't start: ", err)
	}
}
