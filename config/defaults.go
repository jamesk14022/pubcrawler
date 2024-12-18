package config

import (
	types "github.com/jamesk14022/barcrawler/types"
)

// Define the default config
var DefaultConfig = types.MarkerSettings{
	2: {
		"distanceThreshold": 1.5,
		"mu":                1.5,
		"alpha":             1.5,
	},
	3: {
		"distanceThreshold": 1.7,
		"mu":                1.7,
		"alpha":             1.7,
	},
	4: {
		"distanceThreshold": 2.0,
		"mu":                1.7,
		"alpha":             1.7,
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
}
