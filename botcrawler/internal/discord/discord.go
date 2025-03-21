package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lepotekil/CIOS/internal/logger"
	"github.com/lepotekil/CIOS/internal/structs"
)

// DiscordInit initializes a connection to Discord with the bot's token.
// Logs the success or failure using Loggog.
func DiscordInit(config structs.Config) (*discordgo.Session, error) {
	// Create a new Discord session using the token from the config
	dg, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error initializing Discord session: %v", err), nil, config)
		return nil, fmt.Errorf("error discordgo.New(): %w", err)
	}

	// Open the WebSocket connection
	err = dg.Open()
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error opening WebSocket connection: %v", err), nil, config)
		return nil, fmt.Errorf("ws connection error: %w", err)
	}

	// Log success using the logger
	logger.Loggog("Successfully connected to Discord!", nil, config)

	return dg, nil
}
