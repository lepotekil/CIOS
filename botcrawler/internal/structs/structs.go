package structs

// PlayerDetails represents the detailed information of a player returned by the API
type PlayerDetails struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	RegistrationDate    string `json:"registrationDate"`
	LastActivityDate    string `json:"lastActivityDate"`
	FactionLastActivity string `json:"factionLastActivityDate"`
	ActivityTime        int    `json:"activityTime"`
	Rank                string `json:"rank,omitempty"` // Optional field
	Power               int    `json:"power"`
	Role                string `json:"role,omitempty"` // Optional field
	Online              bool   `json:"online"`
	OnlineServer        string `json:"onlineServer,omitempty"` // Optional field
	Faction             struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Icon         string `json:"icon"`
		Description  string `json:"description"`
		CreationDate string `json:"creationDate"`
		FirstDay     string `json:"firstDay"`
		LastDay      string `json:"lastDay"`
		UUID         string `json:"uuid"`
	} `json:"faction,omitempty"` // Optional field
	HeadURL string `json:"headUrl"`
}

// Player represents a player with a name and ID
type Player struct {
	Name string // Name of the player
	ID   string // ID of the player
}

type OnlinePlayers struct {
	Props struct {
		PageProps struct {
			Stats struct {
				OnlinePlayers []Player `json:"onlinePlayers"`
			} `json:"stats"`
		} `json:"pageProps"`
	} `json:"props"`
}

// Config structure to hold configuration data
type Config struct {
	// Discord configuration: Holds the information required for Discord bot integration
	Discord struct {
		Token      string `yaml:"token"`       // Bot token for Discord
		LogChannel string `yaml:"log_channel"` // Channel ID for logging errors and messages
	} `yaml:"discord"`

	// Pactify API configuration: URLs for various Pactify API endpoints
	PactifyAPI struct {
		Players string `yaml:"players"` // URL to get the player ID code
	} `yaml:"pactify_api"`

	// Database configuration: Holds the credentials and connection details for the database
	Database struct {
		IP       string `yaml:"ip"`       // Database IP address
		Port     string `yaml:"port"`     // Database port (default is 3306 for MySQL)
		DBName   string `yaml:"db_name"`  // Name of the database
		Username string `yaml:"username"` // Database username
		Password string `yaml:"password"` // Database password
	} `yaml:"database"`

	// Migrations configuration: Path to the folder containing migration files
	Migrations struct {
		Folder string `yaml:"folder"` // Path to the SQL migrations folder
	} `yaml:"migrations"`

	// ProxyList holds a list of proxies to use for API requests
	ProxyList []string `yaml:"proxy_list"` // List of proxy URLs
}
