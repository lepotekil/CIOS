package utils

import (
	"fmt"
	"os"

	"github.com/lepotekil/CIOS/internal/structs"
	"gopkg.in/yaml.v2"
)

// LoadConfig loads configuration from a YAML file and populates the provided config structure.
// Parameters:
// - filename: Path to the YAML configuration file.
// - config: A pointer to the struct where the configuration will be loaded.
// Returns:
// - error: Returns an error if the file cannot be opened or the YAML decoding fails.
func LoadConfig(filename string, config *structs.Config) error {
	// Attempt to open the configuration file
	file, err := os.Open(filename)
	if err != nil {
		// Return a formatted error message if the file cannot be opened
		return fmt.Errorf("failed to open config file '%s': %v", filename, err)
	}
	// Ensure the file is closed once the function exits
	defer file.Close()

	// Create a new YAML decoder for the file
	decoder := yaml.NewDecoder(file)

	// Decode the YAML file into the provided config structure
	if err := decoder.Decode(config); err != nil {
		// Return a formatted error message if decoding fails
		return fmt.Errorf("failed to decode YAML from config file '%s': %v", filename, err)
	}

	// Successfully loaded configuration
	fmt.Printf("Configuration loaded successfully from '%s'.\n", filename)
	return nil
}
