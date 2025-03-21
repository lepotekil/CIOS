package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lepotekil/CIOS/internal/structs"
)

// Global variables to track the current log file and its date
var currentLogFile *os.File
var currentLogDate string

// Loggog writes log messages to a log file named by the current date and sends them to Discord.
func Loggog(message string, dg *discordgo.Session, config structs.Config) {
	now := time.Now()
	today := now.Format("2006-01-02") // Current date in YYYY-MM-DD format

	// Check if the log folder exists, if not, create it
	logDir := "log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// Create the log folder if it does not exist
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			fmt.Printf("[x] Failed to create log directory: %v\n", err)
			return
		}
	}

	// Clean up old log files (5 days or older)
	cleanupOldLogs(logDir, 5)

	// Check if a new log file is needed
	if currentLogDate != today {
		// Close the previous log file if it exists
		if currentLogFile != nil {
			currentLogFile.Close()
		}

		// Open or create a log file for today's date in the "log" folder
		file, err := os.OpenFile(fmt.Sprintf("%s/%s.log", logDir, today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("[x] Failed to open log file: %v\n", err)
			return
		}

		currentLogFile = file
		currentLogDate = today
	}

	// Format the message with a timestamp (without milliseconds)
	timestamp := now.Format("15:04:05") // HH:MM:SS (no milliseconds)
	formattedMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	// Write the message to the log file
	_, err := currentLogFile.WriteString(formattedMessage)
	if err != nil {
		fmt.Printf("[x] Failed to write to log file: %v\n", err)
		return
	}

	// Send the log message to Discord, if a Discord session is provided
	if dg != nil {
		err = sendMessageDiscordChannel(dg, config.Discord.LogChannel, formattedMessage, config)
		if err != nil {
			fmt.Printf("[x] Failed to send log to Discord: %v\n", err)
		}
	}

	// Print to console for immediate feedback
	fmt.Print(formattedMessage)
}

// sendMessageDiscordChannel sends a message to a specific Discord channel.
// Uses Loggog for logging success and error messages.
func sendMessageDiscordChannel(dg *discordgo.Session, channelID, message string, config structs.Config) error {
	// Send the message to the specified channel
	_, err := dg.ChannelMessageSend(channelID, message)
	if err != nil {
		// Log the error using Loggog
		Loggog(fmt.Sprintf("Error sending message to channel %s: %v", channelID, err), nil, config)
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Log the success using Loggog
	Loggog(fmt.Sprintf("Successfully sent message to channel %s", channelID), nil, config)
	return nil
}

// cleanupOldLogs deletes log files older than the specified number of days.
func cleanupOldLogs(logDir string, maxAgeDays int) {
	now := time.Now()
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("[x] Error accessing file: %v\n", err)
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate the file's age
		fileAge := now.Sub(info.ModTime())
		if fileAge.Hours() >= float64(maxAgeDays*24) {
			// Delete the file if it's older than the max age
			err := os.Remove(path)
			if err != nil {
				fmt.Printf("[x] Failed to delete old log file: %s, error: %v\n", path, err)
			} else {
				fmt.Printf("[✓] Deleted old log file: %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("[x] Error during log cleanup: %v\n", err)
	}
}
