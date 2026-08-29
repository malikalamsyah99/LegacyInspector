package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type AppleModel struct {
	SystemProductName string   `yaml:"SystemProductName"`
	AppleModelCode    []string `yaml:"AppleModelCode"`
	AppleModelYear    []int    `yaml:"AppleModelYear"`

	Specifications struct {
		SystemReportName FlexibleStrings `yaml:"SystemReportName"`
		MarketingName    FlexibleStrings `yaml:"MarketingName"`
	} `yaml:"Specifications"`
}

type FlexibleStrings []string

func (f *FlexibleStrings) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		if value.Value == "" {
			*f = nil
			return nil
		}

		*f = []string{value.Value}
		return nil

	case yaml.SequenceNode:
		var values []string

		for _, node := range value.Content {
			if node.Kind == yaml.ScalarNode && node.Value != "" {
				values = append(values, node.Value)
			}
		}

		*f = values
		return nil

	default:
		*f = nil
		return nil
	}
}

type ModelInfo struct {
	Name      string `json:"name"`
	Model     string `json:"model"`
	ModelYear int    `json:"model_year"`
	Source    string `json:"source"`
}

func main() {
	root := "tools/import_models/OpenCorePkg/AppleModels/DataBase"
	outputFile := "internal/decoder/data/models.json"

	var files []string

	err := filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(
				strings.ToLower(info.Name()),
				".yaml",
			) {
				files = append(files, path)
			}

			return nil
		},
	)

	if err != nil {
		panic(err)
	}

	sort.Strings(files)

	models := make(map[string]ModelInfo)

	processedFiles := 0
	processedCodes := 0
	failedFiles := 0

	for _, file := range files {
		data, err := os.ReadFile(file)

		if err != nil {
			fmt.Printf(
				"WARNING: failed to read %s: %v\n",
				file,
				err,
			)

			failedFiles++
			continue
		}

		var model AppleModel

		if err := yaml.Unmarshal(data, &model); err != nil {
			fmt.Printf(
				"WARNING: failed to parse %s: %v\n",
				file,
				err,
			)

			failedFiles++
			continue
		}

		if model.SystemProductName == "" {
			fmt.Printf(
				"WARNING: %s has no SystemProductName\n",
				file,
			)

			failedFiles++
			continue
		}

		if len(model.AppleModelCode) == 0 {
			fmt.Printf(
				"WARNING: %s has no AppleModelCode\n",
				file,
			)

			failedFiles++
			continue
		}

		modelYear := 0

		if len(model.AppleModelYear) > 0 {
			modelYear = model.AppleModelYear[0]
		}

		name := model.SystemProductName

		if len(model.Specifications.MarketingName) > 0 {
			name = model.Specifications.MarketingName[0]
		} else if len(model.Specifications.SystemReportName) > 0 {
			name = model.Specifications.SystemReportName[0]
		}

		source := filepath.Base(file)

		for _, code := range model.AppleModelCode {
			code = strings.ToUpper(
				strings.TrimSpace(code),
			)

			if code == "" {
				continue
			}

			models[code] = ModelInfo{
				Name:      name,
				Model:     model.SystemProductName,
				ModelYear: modelYear,
				Source:    source,
			}

			processedCodes++
		}

		processedFiles++
	}

	output, err := json.MarshalIndent(
		models,
		"",
		"  ",
	)

	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		outputFile,
		output,
		0644,
	); err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("LegacyInspector Model Importer")
	fmt.Println("================================")
	fmt.Printf(
		"YAML files found     : %d\n",
		len(files),
	)
	fmt.Printf(
		"Files imported       : %d\n",
		processedFiles,
	)
	fmt.Printf(
		"Files failed         : %d\n",
		failedFiles,
	)
	fmt.Printf(
		"Model codes imported : %d\n",
		processedCodes,
	)
	fmt.Printf(
		"Unique model codes   : %d\n",
		len(models),
	)
	fmt.Printf(
		"Output               : %s\n",
		outputFile,
	)
	fmt.Println("================================")
}
