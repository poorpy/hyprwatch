package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type EventConfig struct {
	events map[string][]Command
}

type Command struct {
	Data     string `yaml:"data"`
	Callback string `yaml:"callback"`
}

func NewConfig(configPath string) (EventConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return EventConfig{}, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	eventMap := make(map[string][]Command, 0)
	if err := yaml.Unmarshal(data, &eventMap); err != nil {
		return EventConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return EventConfig{
		events: eventMap,
	}, nil
}

func (wc *EventConfig) Lookup(eventName string) ([]Command, bool) {
	commands, ok := wc.events[eventName]
	return commands, ok
}
