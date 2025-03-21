-- Main table to store player information
CREATE TABLE IF NOT EXISTS botcrawler_players (
    id INT AUTO_INCREMENT PRIMARY KEY, -- Unique identifier for the player in the database
    player_id VARCHAR(255) NOT NULL UNIQUE, -- Unique player ID (fixed and immutable)
    current_player_name VARCHAR(255) NOT NULL, -- Current name of the player
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Date and time when the player was added
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Last update timestamp for the player
) ENGINE=InnoDB;

-- Table to store the history of player names
CREATE TABLE IF NOT EXISTS botcrawler_player_name_history (
    id INT AUTO_INCREMENT PRIMARY KEY, -- Unique identifier for the name history entry
    player_id VARCHAR(255) NOT NULL, -- Reference to the player_id in the botcrawler_players table
    previous_name VARCHAR(255) NOT NULL, -- Previous name of the player
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the name change occurred
    CONSTRAINT fk_botcrawler_player FOREIGN KEY (player_id) REFERENCES botcrawler_players(player_id) ON DELETE CASCADE -- Foreign key linking to botcrawler_players
) ENGINE=InnoDB;

-- Index to optimize searches on player_id in the history table
CREATE INDEX idx_botcrawler_player_name_history_player_id ON botcrawler_player_name_history(player_id);