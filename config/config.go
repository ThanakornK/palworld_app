package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	GinMode        string
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	DataDir        string
	PalsFile       string
	StoredPalsFile string
	PassiveSkillsFile string
	PassiveSkillCombosFile string
}

// LoadConfig loads configuration from environment variables with defaults
// It first attempts to load from .env file, then uses environment variables or defaults
func LoadConfig() *Config {
	// Try to load .env file (ignore errors if file doesn't exist)
	LoadEnvFile(".env")
	
	return &Config{
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "release"),
		AllowedOrigins: getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001"}),
		AllowedMethods: getEnvSlice("ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		AllowedHeaders: getEnvSlice("ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
		DataDir:        getEnv("DATA_DIR", "./data"),
		PalsFile:       getEnv("PALS_FILE", "pals.json"),
		StoredPalsFile: getEnv("STORED_PALS_FILE", "stored_pals.json"),
		PassiveSkillsFile: getEnv("PASSIVE_SKILLS_FILE", "passive_skills.json"),
		PassiveSkillCombosFile: getEnv("PASSIVE_SKILL_COMBOS_FILE", "passive_skill_combos.json"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvSlice gets an environment variable as a slice with a default value
func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}