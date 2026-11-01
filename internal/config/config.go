package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HSLAPIKey  string
	ArrivalIDs []string
	Debug      bool
}

func Load() Config {
	hslAPIKey := os.Getenv("HSL_API_KEY")

	if hslAPIKey == "" {
		log.Fatal("HSL_API_KEY environment variable is required")
	}

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	return Config{
		HSLAPIKey:  hslAPIKey,
		ArrivalIDs: splitIDs(os.Getenv("ARRIVAL_IDS")),
		Debug:      debug,
	}
}

func splitIDs(s string) []string {
	var ids []string
	for part := range strings.SplitSeq(s, ",") {
		if id := strings.TrimSpace(part); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}
