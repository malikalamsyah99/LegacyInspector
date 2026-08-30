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

//go:embed data/device_identifiers.json
var deviceIdentifierData []byte

var models map[string]ModelInfo
var modelOverrides map[string]ModelInfo
var deviceIdentifiers map[string]string

func init() {
	models = make(map[string]ModelInfo)
	modelOverrides = make(map[string]ModelInfo)
	deviceIdentifiers = make(map[string]string)

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

	if err := json.Unmarshal(
		deviceIdentifierData,
		&deviceIdentifiers,
	); err != nil {
		panic(fmt.Sprintf(
			"failed to load device identifier database: %v",
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

func lookupModelByDeviceName(name string) (ModelInfo, bool) {
	model, ok := deviceIdentifiers[name]
	if !ok {
		return ModelInfo{}, false
	}

	return ModelInfo{
		Name:   name,
		Model:  model,
		Source: "apple_device_identifiers",
	}, true
}
