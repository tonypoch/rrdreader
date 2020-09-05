package rrdtool

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

type Config struct {
	MultipleSources bool
	Directory       string
	VLabel          string
	BadPrefixes     []string
}

var Configs = map[string]Config{
	"Network": {
		MultipleSources: false,
		Directory:       "interface-e*",
		VLabel:          "Packets in eth interface",
		BadPrefixes:     []string{"if_"},
	},
	"Memory": {
		MultipleSources: true,
		Directory:       "memory*",
		VLabel:          "Memory usage",
		BadPrefixes:     []string{"memory-"},
	},
	"Disk": {
		MultipleSources: false,
		Directory:       "disk-sd?",
		VLabel:          "Disk operations",
		BadPrefixes:     []string{"disk_"},
	},
	"CPU": {
		MultipleSources: true,
		Directory:       "cpu-*",
		VLabel:          "CPU usage",
		BadPrefixes:     []string{"cpu-"},
	},
	"Load": {
		MultipleSources: false,
		Directory:       "load*",
		VLabel:          "Load on system",
	},
	"Temp": {
		MultipleSources: true,
		Directory:       "sensors-*temp*",
		VLabel:          "Sensors temperature",
		BadPrefixes:     []string{"temperature-"},
	},
	"SMART": {
		MultipleSources: false,
		Directory:       "smart-sd?",
		VLabel:          "S.M.A.R.T Attributes",
		BadPrefixes:     []string{"smart_", "attribute-"},
	},
}

var GraphPallete = []string{
	//X11 color names
	"FF0000", //red
	"0000FF", //blue
	"008000", //green
	"FFC0CB", //pink
	"FFA500", //orange
	"8A2BE2", //blue violet
	"008080", //teal
	"808080", //grey
	"006400", //dark green
	"ADD8E6", //lightblue
	"A52A2A", //brown
	"800080", //purple
	"FF6347", //tomato
	"F5DEB3", //wheat
}

func GenerateHexColor() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetTimePeriods(now time.Time) map[string]time.Time {
	return map[string]time.Time{
		"day":   now.AddDate(0, 0, -1),
		"week":  now.AddDate(0, 0, -7),
		"month": now.AddDate(0, -1, 0),
		"year":  now.AddDate(-1, 0, 0),
	}
}

func RemovePrefix(filename string, prefixes []string) string {
	result := filename
	for _, prefix := range prefixes {
		result = strings.TrimPrefix(result, prefix)
	}
	return result
}
