package decoder

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/locations.json
var locationData []byte

var locations map[string]Location

func init() {
	locations = make(map[string]Location)

	if err := json.Unmarshal(locationData, &locations); err != nil {
		panic(fmt.Sprintf(
			"failed to load location database: %v",
			err,
		))
	}
}

func lookupLocation(code string, format string) (Location, bool) {
	location, ok := locations[code]

	if !ok {
		return Location{}, false
	}

	if location.Format != format {
		return Location{}, false
	}

	return location, true
}
