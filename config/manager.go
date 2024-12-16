package config

import (
	types "github.com/jamesk14022/barcrawler/types"
)

// ConfigManager manages multiple versions of marker settings.
type ConfigManager struct {
	ConfigVersions map[string]types.MarkerSettings
}

// NewConfigManager creates a new ConfigManager with predefined versions.
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		ConfigVersions: map[string]types.MarkerSettings{
			"tallinn": {
				3: {
					"distanceThreshold": 1.0,
					"mu":                1.5,
					"alpha":             1.5,
				},
				4: {
					"distanceThreshold": 1.0,
					"mu":                1.1,
					"alpha":             1.1,
				},
				5: {
					"distanceThreshold": 1.6,
					"mu":                1.3,
					"alpha":             1.3,
				},
				6: {
					"distanceThreshold": 1.6,
					"mu":                1.3,
					"alpha":             1.3,
				},
			},
			"edinburgh": {
				2: {
					"distanceThreshold": 1.5,
					"mu":                1.5,
					"alpha":             1.5,
				},
				3: {
					"distanceThreshold": 1.5,
					"mu":                1.5,
					"alpha":             1.5,
				},
				4: {
					"distanceThreshold": 1.5,
					"mu":                1.5,
					"alpha":             1.5,
				},
				5: {
					"distanceThreshold": 2.0,
					"mu":                1.7,
					"alpha":             1.7,
				},
				6: {
					"distanceThreshold": 2.0,
					"mu":                1.7,
					"alpha":             1.7,
				},
			},
		},
	}
}

// GetConfig retrieves the configuration for a specific city or returns the default configuration.
func (cm *ConfigManager) GetConfig(city string) types.MarkerSettings {
	if config, exists := cm.ConfigVersions[city]; exists {
		return config
	}
	return DefaultConfig
}
