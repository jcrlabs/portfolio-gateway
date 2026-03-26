package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         string
	RoutingFile  string
	AllowOrigins string
	RateLimitRPS float64
}

type Route struct {
	Name    string `yaml:"name"`
	Prefix  string `yaml:"prefix"`
	Backend string `yaml:"backend"`
}

type routingFile struct {
	Routes []Route `yaml:"routes"`
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Port:         getEnv("PORT", "8080"),
		RoutingFile:  getEnv("ROUTING_FILE", "routing.yaml"),
		AllowOrigins: getEnv("ALLOW_ORIGINS", "https://jcrlabs.net"),
		RateLimitRPS: 100,
	}
}

func LoadRoutes(path string) []Route {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read routing file %q: %v", path, err)
	}
	var rf routingFile
	if err = yaml.Unmarshal(data, &rf); err != nil {
		log.Fatalf("parse routing file: %v", err)
	}
	return rf.Routes
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
