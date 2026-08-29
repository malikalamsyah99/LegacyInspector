package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Location struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Format string `json:"format"`
}

func parseSection(filename, marker string) []string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	var values []string
	inSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, marker) {
			inSection = true
			continue
		}

		if !inSection {
			continue
		}

		if strings.HasPrefix(line, "};") {
			break
		}

		if !strings.HasPrefix(line, "\"") {
			continue
		}

		end := strings.Index(line[1:], "\"")

		if end < 0 {
			continue
		}

		value := line[1 : end+1]
		values = append(values, value)
	}

	return values
}

func main() {
	inputFile := "tools/modelinfo.h"
	outputFile := "internal/decoder/data/locations.json"

	legacyCodes := parseSection(
		inputFile,
		"static const char *AppleLegacyLocations[]",
	)

	legacyNames := parseSection(
		inputFile,
		"static const char *AppleLegacyLocationNames[]",
	)

	modernCodes := parseSection(
		inputFile,
		"static const char *AppleLocations[]",
	)

	modernNames := parseSection(
		inputFile,
		"static const char *AppleLocationNames[]",
	)

	if len(legacyCodes) != len(legacyNames) {
		panic(fmt.Sprintf(
			"legacy location mismatch: %d codes, %d names",
			len(legacyCodes),
			len(legacyNames),
		))
	}

	if len(modernCodes) != len(modernNames) {
		panic(fmt.Sprintf(
			"modern location mismatch: %d codes, %d names",
			len(modernCodes),
			len(modernNames),
		))
	}

	locations := make(map[string]Location)

	for i := range legacyCodes {
		locations[legacyCodes[i]] = Location{
			Code:   legacyCodes[i],
			Name:   legacyNames[i],
			Format: "legacy-11",
		}
	}

	for i := range modernCodes {
		locations[modernCodes[i]] = Location{
			Code:   modernCodes[i],
			Name:   modernNames[i],
			Format: "legacy-12",
		}
	}

	data, err := json.MarshalIndent(
		locations,
		"",
		"  ",
	)

	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		outputFile,
		data,
		0644,
	); err != nil {
		panic(err)
	}

	fmt.Printf(
		"Imported %d legacy-11 locations\n",
		len(legacyCodes),
	)

	fmt.Printf(
		"Imported %d legacy-12 locations\n",
		len(modernCodes),
	)

	fmt.Printf(
		"Total: %d locations\n",
		len(locations),
	)

	fmt.Printf(
		"Output: %s\n",
		outputFile,
	)
}
