package decoder

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/products.json
var productData []byte

//go:embed data/macserial_products.json
var macSerialProductData []byte

//go:embed data/overrides.json
var overrideProductData []byte

var products map[string]Product
var macSerialProducts map[string]Product
var overrideProducts map[string]Product

func init() {
	products = make(map[string]Product)
	macSerialProducts = make(map[string]Product)
	overrideProducts = make(map[string]Product)

	if err := json.Unmarshal(productData, &products); err != nil {
		panic(fmt.Sprintf(
			"failed to load product database: %v",
			err,
		))
	}

	if err := json.Unmarshal(
		macSerialProductData,
		&macSerialProducts,
	); err != nil {
		panic(fmt.Sprintf(
			"failed to load MacSerials product database: %v",
			err,
		))
	}

	if err := json.Unmarshal(
		overrideProductData,
		&overrideProducts,
	); err != nil {
		panic(fmt.Sprintf(
			"failed to load product override database: %v",
			err,
		))
	}
}

func lookupProduct(code string) (Product, bool) {
	if product, ok := overrideProducts[code]; ok {
		return product, true
	}

	if product, ok := products[code]; ok {
		return product, true
	}

	if product, ok := macSerialProducts[code]; ok {
		return product, true
	}

	return Product{}, false
}
