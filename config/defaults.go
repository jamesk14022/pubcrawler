package config

import (
	types "github.com/jamesk14022/barcrawler/types"
)

// Define the default config
var DefaultConfig = types.MarkerSettings{
	3: {
		"distanceThreshold": 0.9,
		"mu":                1.1,
		"alpha":             1.1,
	},
	4: {
		"distanceThreshold": 0.9,
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
}
