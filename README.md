Welcome to my project!

## Database Structure -- 
![Database setup](web/DBstruct.svg)
# LoL International Rankings API

A polling service for League of Legends casters, journalists, and analysts to create collaborative top 25 international team rankings.

## 🎯 Overview

With international League of Legends matches being few and far between, it's challenging to understand how teams from different regions stack up against each other. This service provides a platform for qualified voters to submit weekly rankings, creating a comprehensive view of the global competitive landscape.

# LoL International Rankings API

A polling service for League of Legends casters, journalists, and analysts to create collaborative top 25 international team rankings.

## Tech Stack

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL
- **Port**: 31022

## API Endpoints

Base URL: `http://localhost:31022/api`

### Health Check
- `GET /health` - Check if the service is running

### Teams
- `GET /teams` - Retrieve all teams available for ranking

### Voters
- `GET /voters` - Get list of approved voters
- `POST /voters` - Apply to become a voter or add a new voter

### Rankings
- `POST /rankings` - Submit weekly rankings (authenticated voters only)
- `GET /rankings` - Get current week's aggregated rankings
- `GET /rankings/week/:week` - Get rankings for a specific week

## Database Schema

### poll_periods
- `id` (integer, primary key)
- `name` (varchar 50, not null)
- `start_date` (date, not null)
- `end_date` (date, not null)
- `is_active` (boolean, default true)
- `created_at` (timestamp, default CURRENT_TIMESTAMP)

### teams
- `id` (integer, primary key)
- `name` (varchar 100, not null)
- `region` (varchar 10, not null)
- `logo_url` (varchar 500)
- `created_at` (timestamp, default CURRENT_TIMESTAMP)

### voters
- `id` (integer, primary key)
- `name` (varchar 100, not null)
- `email` (varchar 200, not null, unique)
- `role` (varchar 50, default 'voter')
- `active` (boolean, default true)
- `created_at` (timestamp, default CURRENT_TIMESTAMP)

### votes
- `id` (integer, primary key)
- `voter_id` (integer, foreign key to voters.id)
- `team_id` (integer, foreign key to teams.id)
- `poll_period_id` (integer, foreign key to poll_periods.id)
- `rank_position` (integer, not null)
- `created_at` (timestamp, default CURRENT_TIMESTAMP)
- Unique constraints on (voter_id, rank_position, poll_period_id) and (voter_id, team_id, poll_period_id)

## Roadmap

- JWT Authentication
- Admin panel
- Voter approval workflow
- Email notifications
- Historical rankings API
- Public frontend
- Webhook integrations
- Rate limiting
- API versioning
- Comprehensive test suite

## 📞 Contact

- Project maintainer: Lukas Dutrieux
- Email: dutrieuxl31022@gmail.com
---
