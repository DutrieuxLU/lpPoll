// cmd/database.go
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DatabaseService struct {
	db     *sql.DB
	teams  []Team
	voters []Voter
	votes  []Vote
}

func NewDatabaseService(db *sql.DB) *DatabaseService {
	return &DatabaseService{
		db:     db,
		teams:  []Team{},
		voters: []Voter{},
		votes:  []Vote{},
	}
}

func (ds *DatabaseService) LoadFromDB() error {
	// Load teams
	if err := ds.loadTeams(); err != nil {
		return err
	}

	// Load voters
	if err := ds.loadVoters(); err != nil {
		return err
	}

	// Load votes
	if err := ds.loadVotes(); err != nil {
		return err
	}

	log.Printf("Loaded %d teams, %d voters, and %d votes from database",
		len(ds.teams), len(ds.voters), len(ds.votes))
	return nil
}

func (ds *DatabaseService) loadTeams() error {
	rows, err := ds.db.Query("SELECT id, name, region FROM teams ORDER BY name")
	if err != nil {
		return fmt.Errorf("failed to query teams: %v", err)
	}
	defer rows.Close()

	ds.teams = []Team{} // Clear existing teams
	for rows.Next() {
		var team Team
		if err = rows.Scan(&team.ID, &team.Name, &team.Region); err != nil {
			return fmt.Errorf("failed to scan team: %v", err)
		}
		ds.teams = append(ds.teams, team)
	}

	return rows.Err()
}

func (ds *DatabaseService) loadVoters() error {
	rows, err := ds.db.Query("SELECT id, name, email, role, active, created_at FROM voters ORDER BY name")
	if err != nil {
		return fmt.Errorf("failed to query voters: %v", err)
	}
	defer rows.Close()

	ds.voters = []Voter{} // Clear existing voters
	for rows.Next() {
		var voter Voter
		if err = rows.Scan(&voter.ID, &voter.Name, &voter.Email, &voter.Role, &voter.Active, &voter.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan voter: %v", err)
		}
		ds.voters = append(ds.voters, voter)
	}

	return rows.Err()
}

func (ds *DatabaseService) loadVotes() error {
	rows, err := ds.db.Query("SELECT id, voter_id, team_id, poll_period_id, rank_position, created_at FROM votes ORDER BY poll_period_id DESC, voter_id, rank_position")
	if err != nil {
		return fmt.Errorf("failed to query votes: %v", err)
	}
	defer rows.Close()

	ds.votes = []Vote{} // Clear existing votes
	for rows.Next() {
		var vote Vote
		if err = rows.Scan(&vote.ID, &vote.VoterID, &vote.TeamID, &vote.PollPeriodID, &vote.RankPosition, &vote.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan vote: %v", err)
		}
		ds.votes = append(ds.votes, vote)
	}

	return rows.Err()
}

func (ds *DatabaseService) GetTeams() []Team {
	return ds.teams
}

func (ds *DatabaseService) GetVoters() []Voter {
	return ds.voters
}

func (ds *DatabaseService) GetVotes() []Vote {
	return ds.votes
}

func (ds *DatabaseService) CreateVoter(name, email, role string) (*Voter, error) {
	// Since role has a default value and created_at has CURRENT_TIMESTAMP default,
	// we can let the database handle those if not specified
	query := `INSERT INTO voters (name, email, role) 
			  VALUES ($1, $2, $3) RETURNING id, role, active, created_at`

	var voter Voter
	voter.Name = name
	voter.Email = email

	err := ds.db.QueryRow(query, name, email, role).Scan(&voter.ID, &voter.Role, &voter.Active, &voter.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create voter: %v", err)
	}

	// Add to in-memory cache
	ds.voters = append(ds.voters, voter)

	return &voter, nil
}

func (ds *DatabaseService) SubmitRanking(voterID, pollPeriodID int, teams []int) error {
	// Start transaction
	tx, err := ds.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Delete existing votes for this voter and poll period
	_, err = tx.Exec("DELETE FROM votes WHERE voter_id = $1 AND poll_period_id = $2", voterID, pollPeriodID)
	if err != nil {
		return fmt.Errorf("failed to delete existing votes: %v", err)
	}

	// Insert new votes
	for rankPosition, teamID := range teams {
		_, err = tx.Exec(`INSERT INTO votes (voter_id, team_id, poll_period_id, rank_position, created_at) 
						  VALUES ($1, $2, $3, $4, $5)`,
			voterID, teamID, pollPeriodID, rankPosition+1, time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert vote: %v", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Reload votes from database to update in-memory cache
	return ds.loadVotes()
}

func (ds *DatabaseService) CalculateRankings(pollPeriodID int) []RankingResponse {
	teamPoints := make(map[int]int)
	teamFirstVotes := make(map[int]int)

	for _, vote := range ds.votes {
		if vote.PollPeriodID == pollPeriodID {
			points := max(0, 26-vote.RankPosition)
			teamPoints[vote.TeamID] += points

			if vote.RankPosition == 1 {
				teamFirstVotes[vote.TeamID]++
			}
		}
	}

	// Convert to response format and sort
	var rankings []RankingResponse
	for _, team := range ds.teams {
		if points, exists := teamPoints[team.ID]; exists {
			rankings = append(rankings, RankingResponse{
				TeamID:     team.ID,
				TeamName:   team.Name,
				Points:     points,
				FirstVotes: teamFirstVotes[team.ID],
			})
		}
	}

	// Sort by points (bubble sort for simplicity)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
