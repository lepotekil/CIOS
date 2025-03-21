package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lepotekil/CIOS/internal/database"
	"github.com/lepotekil/CIOS/internal/discord"
	"github.com/lepotekil/CIOS/internal/logger"
	"github.com/lepotekil/CIOS/internal/migrations"
	"github.com/lepotekil/CIOS/internal/structs"
	"github.com/lepotekil/CIOS/internal/utils"
	"golang.org/x/exp/rand"
	"golang.org/x/net/html"
)

func processPlayer(player structs.Player, wg *sync.WaitGroup, dg *discordgo.Session, db *sql.DB) {
	defer wg.Done()

	// Check if the player exists in the database
	exists, err := database.CheckIfPlayerExists(db, player.ID)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error checking if player exists (ID: %s): %v", player.ID, err), dg, config)
		return
	}

	if !exists {
		// If player doesn't exist, add them to the database
		err = database.AddPlayer(db, player.ID, player.Name)
		if err != nil {
			logger.Loggog(fmt.Sprintf("Error adding new player (ID: %s, Name: %s): %v", player.ID, player.Name, err), dg, config)
			return
		}
		return
	}

	// If the player exists, check if the name has changed
	var currentName string
	query := "SELECT current_player_name FROM botcrawler_players WHERE player_id = ?"
	err = db.QueryRow(query, player.ID).Scan(&currentName)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error retrieving current player name (ID: %s): %v", player.ID, err), dg, config)
		return
	}

	if currentName != player.Name {
		// Update the current name and add the old name to the history
		err = database.UpdatePlayerName(db, player.ID, player.Name)
		if err != nil {
			logger.Loggog(fmt.Sprintf("Error updating player name (ID: %s): %v", player.ID, err), dg, config)
			return
		}

		err = database.AddNameToHistory(db, player.ID, currentName)
		if err != nil {
			logger.Loggog(fmt.Sprintf("Error adding name to history (ID: %s, Previous Name: %s): %v", player.ID, currentName, err), dg, config)
			return
		}

		// Retrieve and log the player's name history
		history, err := database.GetNameHistory(db, player.ID)
		if err != nil {
			logger.Loggog(fmt.Sprintf("Error retrieving name history (ID: %s): %v", player.ID, err), dg, config)
			return
		}

		// Log the player's name change in French (short message)
		logger.Loggog(fmt.Sprintf("Joueur '%s' -> '%s', historique : %v", currentName, player.Name, history), dg, config)
	}
}

// getRandomProxy selects a random proxy from the list
func getRandomProxy(proxyList []string) string {
	randomIndex := rand.Intn(len(proxyList))
	return proxyList[randomIndex]
}

// fetchPlayers retrieves the list of players from the given URL using a random proxy.
func fetchPlayers(urlStr string, config structs.Config, dg *discordgo.Session) ([]structs.Player, error) {

	// Select a random proxy from the list
	selectedProxy := getRandomProxy(config.ProxyList)

	// Set up the HTTP client with the selected proxy
	proxyURL, err := url.Parse(selectedProxy)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Invalid proxy URL %s: %v", selectedProxy, err), dg, config)
		return nil, err
	}
	httpTransport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: httpTransport}

	// Make the HTTP request using the proxy
	resp, err := client.Get(urlStr)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error fetching data from URL %s using proxy %s: %v", urlStr, selectedProxy, err), dg, config)
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML content
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error parsing HTML from URL %s: %v", urlStr, err), dg, config)
		return nil, err
	}

	// Extract JSON data from the <script> tag
	var jsonData string
	var findJSON func(*html.Node)
	findJSON = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "script" {
			for _, attr := range node.Attr {
				if attr.Key == "id" && attr.Val == "__NEXT_DATA__" {
					jsonData = node.FirstChild.Data
					return
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findJSON(child)
		}
	}
	findJSON(doc)

	// Check if JSON data was found
	if jsonData == "" {
		err := fmt.Errorf("JSON data not found at URL %s", urlStr)
		logger.Loggog(err.Error(), dg, config)
		return nil, err
	}

	// Parse JSON data into the OnlinePlayers structure
	var data structs.OnlinePlayers
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error decoding JSON data from URL %s: %v", urlStr, err), dg, config)
		return nil, err
	}

	// Return the list of players
	return data.Props.PageProps.Stats.OnlinePlayers, nil
}

// Global variable to hold configuration
var config structs.Config

func main() {
	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))

	// Load the configuration from a JSON file
	err := utils.LoadConfig("config/config.yaml", &config)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error loading config: %v", err), nil, config)
		panic(err)
	}

	// Check if the proxy list is empty
	if len(config.ProxyList) == 0 {
		logger.Loggog("Proxy list is empty. Cannot proceed with the request.", nil, config)
		os.Exit(1)
	}

	// Initialize the Discord session using the Discord token
	dg, err := discord.DiscordInit(config)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error initializing Discord: %v", err), nil, config)
		panic(err)
	}
	defer dg.Close()

	// Initialize the database connection using the loaded configuration
	err = database.InitializeDB(config, dg)
	if err != nil {
		logger.Loggog(fmt.Sprintf("Error initializing database: %v", err), dg, config)
		panic(err)
	}
	defer database.CloseDB(dg, config)

	time.Sleep(1 * time.Second)

	// Run the migrations from the folder specified in the configuration
	migrationFolder := config.Migrations.Folder
	err = migrations.RunMigrations(database.DB, migrationFolder, dg, config)
	if err != nil {
		logger.Loggog(fmt.Sprintf("SQL migration error: %v", err), dg, config)
		panic(err)
	}

	time.Sleep(1 * time.Second)

	for {
		// Call the fetchPlayers function
		players, err := fetchPlayers(config.PactifyAPI.Players, config, dg)
		if err != nil {
			// Error is already logged via Loggog, you can handle additional actions if needed
			time.Sleep(10 * time.Second) // Wait before retrying
			continue
		}

		// Create a WaitGroup for concurrent player processing
		var wg sync.WaitGroup

		// Process the list of players
		for _, player := range players {
			// Convert player name to lowercase
			player.Name = strings.ToLower(player.Name)

			// Increment WaitGroup counter
			wg.Add(1)

			// Process the player in a separate goroutine
			go processPlayer(player, &wg, dg, database.DB)
		}

		// Wait for all player processing goroutines to finish
		wg.Wait()

		// Wait before the next fetch
		time.Sleep(120 * time.Second) // Adjust the interval as needed
	}
}
