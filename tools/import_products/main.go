package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Product struct {
	Name string `json:"name"`
}

func main() {

	inputFile := "tools/modelinfo_autogen.h"
	outputFile := "internal/decoder/data/products.json"

	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	products := make(map[string]Product)

	scanner := bufio.NewScanner(file)

	inProductDesc := false

	re := regexp.MustCompile(
		`^\{"([A-Z0-9]{3,4})",\s*"([^"]*)"\},`,
	)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		// Start of AppleModelDesc section
		if strings.Contains(
			line,
			"static APPLE_MODEL_DESC AppleModelDesc[]",
		) {
			inProductDesc = true
			continue
		}

		if !inProductDesc {
			continue
		}

		// End of AppleModelDesc section
		if strings.HasPrefix(line, "};") {
			break
		}

		matches := re.FindStringSubmatch(line)

		if len(matches) != 3 {
			continue
		}

		code := matches[1]
		name := matches[2]

		if name == "" {
			continue
		}

		products[code] = Product{
			Name: name,
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
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
		"Imported %d product codes\n",
		len(products),
	)

	fmt.Printf(
		"Output: %s\n",
		outputFile,
	)
}
