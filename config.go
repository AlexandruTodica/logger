package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type config struct {
	Level  string `json:"level"`
	Output string `json:"output"`
	Parser string `json:"parser"`
}

const (
	configFileEnv = "LOG_CONFIG"
)

func loadConfig() (*config, error) {
	confFile := strings.TrimSpace(os.Getenv(configFileEnv))
	if confFile == "" {
		confFile = "config.json"
	}
	body, err := os.ReadFile(confFile)
	if err != nil {
		return nil, fmt.Errorf("error at reading config file: %w", err)
	}

	var c config
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, fmt.Errorf("error at parsing config file: %w", err)
	}
	return &c, nil
}
