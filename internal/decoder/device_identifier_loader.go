package decoder

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadDeviceIdentifiers(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read device identifiers: %w", err)
	}

	var identifiers map[string]string
	if err := json.Unmarshal(data, &identifiers); err != nil {
		return nil, fmt.Errorf("parse device identifiers: %w", err)
	}

	return identifiers, nil
}
