package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Product struct {
	Name string `json:"name"`
}

var entryPattern = regexp.MustCompile(
	`(?m)^([A-Z0-9]{4}): .*?<configCode>(.*?)</configCode>`,
)

func main() {
	inputFile := "tools/import_macserials/MacModels-D000-DZZZ.xml"
	outputFile := "internal/decoder/data/macserial_products.json"

	data, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	matches := entryPattern.FindAllStringSubmatch(
		string(data),
		-1,
	)

	products := make(map[string]Product)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		code := strings.TrimSpace(match[1])
		name := strings.TrimSpace(match[2])

		if code == "" || name == "" {
			continue
		}

		products[code] = Product{
			Name: name,
		}
	}

	output, err := json.MarshalIndent(
		products,
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

	fmt.Printf(
		"Imported %d MacSerials product codes\n",
		len(products),
	)

	fmt.Printf(
		"Output: %s\n",
		outputFile,
	)
}
