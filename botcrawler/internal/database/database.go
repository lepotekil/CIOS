package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Importing the MySQL driver
	"github.com/lepotekil/CIOS/internal/logger"
	"github.com/lepotekil/CIOS/internal/structs"

	"github.com/bwmarrin/discordgo"
)

// Global variable to hold the database connection
var DB *sql.DB

// InitializeDB initializes a connection to the database using configuration parameters.
func InitializeDB(config structs.Config, dg *discordgo.Session) error {
	// Form the connection string using values from the config
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.IP,
		config.Database.Port,
		config.Database.DBName)

	// Open a connection to the database
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		// Log the error and return
		logger.Loggog(fmt.Sprintf("Error opening database: %v", err), dg, config)
		return fmt.Errorf("error opening database: %v", err)
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		// Log the error and return
		logger.Loggog(fmt.Sprintf("Error pinging database: %v", err), dg, config)
		return fmt.Errorf("error pinging database: %v", err)
	}

	// Set the global DB variable to the opened connection
	DB = db

	// Log the success message
	logger.Loggog("Successfully connected to the database!", dg, config)

	return nil
}

// CloseDB closes the database connection
func CloseDB(dg *discordgo.Session, config structs.Config) {
	if DB != nil {
		// Close the connection and handle any error
		err := DB.Close()
		if err != nil {
			// Log error when closing the connection fails
			logger.Loggog(fmt.Sprintf("Error closing database connection: %v", err), dg, config)
		} else {
			// Log success message when closing the connection
			logger.Loggog("Database connection closed.", dg, config)
		}
	}
}

// CheckIfPlayerExists checks if a player exists in the botcrawler_players table by player_id.
func CheckIfPlayerExists(db *sql.DB, playerID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM botcrawler_players WHERE player_id = ?)"
	err := db.QueryRow(query, playerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if player exists: %w", err)
	}
	return exists, nil
}

// AddPlayer inserts a new player into the botcrawler_players table.
func AddPlayer(db *sql.DB, playerID, playerName string) error {
	query := "INSERT INTO botcrawler_players (player_id, current_player_name) VALUES (?, ?)"
	_, err := db.Exec(query, playerID, playerName)
	if err != nil {
		return fmt.Errorf("failed to add player: %w", err)
	}
	return nil
}

// UpdatePlayerName updates the current_player_name of an existing player.
func UpdatePlayerName(db *sql.DB, playerID, newName string) error {
	query := "UPDATE botcrawler_players SET current_player_name = ? WHERE player_id = ?"
	_, err := db.Exec(query, newName, playerID)
	if err != nil {
		return fmt.Errorf("failed to update player name: %w", err)
	}
	return nil
}

// AddNameToHistory adds an entry into the botcrawler_player_name_history table.
func AddNameToHistory(db *sql.DB, playerID, previousName string) error {
	query := "INSERT INTO botcrawler_player_name_history (player_id, previous_name) VALUES (?, ?)"
	_, err := db.Exec(query, playerID, previousName)
	if err != nil {
		return fmt.Errorf("failed to add name to history: %w", err)
	}
	return nil
}

// GetNameHistory retrieves the name history for a player from botcrawler_player_name_history.
func GetNameHistory(db *sql.DB, playerID string) ([]string, error) {
	query := "SELECT previous_name FROM botcrawler_player_name_history WHERE player_id = ?"
	rows, err := db.Query(query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get name history: %w", err)
	}
	defer rows.Close()

	var history []string
	for rows.Next() {
		var previousName string
		if err := rows.Scan(&previousName); err != nil {
			return nil, fmt.Errorf("failed to scan name history row: %w", err)
		}
		history = append(history, previousName)
	}
	return history, nil
}
