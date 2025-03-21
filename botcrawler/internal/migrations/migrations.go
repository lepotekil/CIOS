package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lepotekil/CIOS/internal/logger"
	"github.com/lepotekil/CIOS/internal/structs"
)

// RunMigrations executes SQL migrations from all files in a directory in sequential order.
func RunMigrations(db *sql.DB, folderPath string, dg *discordgo.Session, config structs.Config) error {
	// Open the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Failed to read the migrations directory: %v", err), dg, config)
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort files by numeric prefix (e.g., 001, 002)
	var sqlFiles []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		sqlFiles = append(sqlFiles, file.Name())
	}
	sort.Strings(sqlFiles)

	// Execute each SQL file in order
	for _, fileName := range sqlFiles {
		filePath := fmt.Sprintf("%s/%s", folderPath, fileName)
		logger.Loggog(fmt.Sprintf("Executing migration file: %s", filePath), dg, config)

		// Read the SQL file content
		query, err := os.ReadFile(filePath)
		if err != nil {
			logger.Loggog(fmt.Sprintf("Failed to read the SQL file: %v", err), dg, config)
			return fmt.Errorf("failed to read SQL file %s: %w", filePath, err)
		}

		// Split the file into individual SQL statements
		queries := strings.Split(string(query), ";")
		for i, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}

			// Execute the SQL statement
			_, execErr := db.Exec(q)
			if execErr != nil {
				logger.Loggog(fmt.Sprintf("Error executing SQL command #%d in file %s: %v\nCommand: %s", i+1, filePath, execErr, q), dg, config)
				return fmt.Errorf("error executing SQL command in file %s: %w", filePath, execErr)
			}
		}
	}

	logger.Loggog("All migrations executed successfully.", dg, config)
	return nil
}
