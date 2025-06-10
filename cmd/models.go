package cmd

import "time"

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

type Voter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
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

type CreateVoterRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}
