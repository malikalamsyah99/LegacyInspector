package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type DeviceIdentifiers map[string]string

func main() {
	root := "tools/import_device_identifiers/apple_device_identifiers"
	input := filepath.Join(root, "devices.json")
	output := "internal/decoder/data/device_identifiers.json"

	data, err := os.ReadFile(input)
	if err != nil {
		panic(fmt.Sprintf("failed to read %s: %v", input, err))
	}

	var devices DeviceIdentifiers
	if err := json.Unmarshal(data, &devices); err != nil {
		panic(fmt.Sprintf("failed to parse %s: %v", input, err))
	}

	keys := make([]string, 0, len(devices))
	for name := range devices {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	result := make(map[string]string, len(devices))
	for _, name := range keys {
		result[name] = devices[name]
	}

	outputData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("failed to encode output: %v", err))
	}

	outputData = append(outputData, '\n')

	if err := os.WriteFile(output, outputData, 0644); err != nil {
		panic(fmt.Sprintf("failed to write %s: %v", output, err))
	}

	fmt.Println("LegacyInspector Device Identifier Importer")
	fmt.Println("================================")
	fmt.Printf("Device identifiers : %d\n", len(result))
	fmt.Printf("Output             : %s\n", output)
	fmt.Println("================================")
}
