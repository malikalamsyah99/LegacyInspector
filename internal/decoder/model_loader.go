package decoder

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type ModelInfo struct {
	Name      string `json:"name"`
	Model     string `json:"model"`
	ModelYear int    `json:"model_year"`
	Source    string `json:"source,omitempty"`
}

//go:embed data/models.json
var modelData []byte

//go:embed data/model_overrides.json
var modelOverrideData []byte

var models map[string]ModelInfo
var modelOverrides map[string]ModelInfo

func init() {
	models = make(map[string]ModelInfo)
	modelOverrides = make(map[string]ModelInfo)

	if err := json.Unmarshal(modelData, &models); err != nil {
		panic(fmt.Sprintf(
			"failed to load model database: %v",
			err,
		))
	}

	if err := json.Unmarshal(
		modelOverrideData,
		&modelOverrides,
	); err != nil {
		panic(fmt.Sprintf(
			"failed to load model override database: %v",
			err,
		))
	}
}

func lookupModel(code string) (ModelInfo, bool) {
	// Manual verified overrides have highest priority.
	if model, ok := modelOverrides[code]; ok {
		return model, true
	}

	// Generated OpenCore model database.
	if model, ok := models[code]; ok {
		return model, true
	}

	return ModelInfo{}, false
}
