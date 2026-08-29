package main

import (
	"fmt"
	"os"

	"legacyinspector/internal/decoder"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("LegacyInspector")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  legacyinspector SERIAL")
		os.Exit(1)
	}

	result := decoder.Decode(os.Args[1])

	fmt.Println()
	fmt.Println("LegacyInspector")
	fmt.Println("================================")
	fmt.Printf("Serial         : %s\n", result.Serial)
	fmt.Printf("Format         : %s\n", result.Format)
	fmt.Printf("Valid          : %t\n", result.Valid)

	if result.Location != "" {
		fmt.Printf("Location Code  : %s\n", result.Location)

		if result.LocationInfo != nil {
			fmt.Printf("Location       : %s\n", result.LocationInfo.Name)
		} else {
			fmt.Println("Location       : Unknown")
		}
	}

	if result.YearCode != "" {
		fmt.Printf("Year Code      : %s\n", result.YearCode)
	}

	if result.Year != 0 {
		fmt.Printf("Production Year : %d\n", result.Year)
	}

	if result.Half != 0 {
		fmt.Printf("Half           : %d\n", result.Half)
	}

	if result.WeekCode != "" {
		fmt.Printf("Week Code      : %s\n", result.WeekCode)
	}

	if result.Week != 0 {
		fmt.Printf("Week           : %d\n", result.Week)
	}

	if result.Production != "" {
		fmt.Printf("Production     : %s\n", result.Production)
	}

	if result.ProductCode != "" {
		fmt.Printf("Product Code   : %s\n", result.ProductCode)
	}

	if result.Product != nil {
		fmt.Printf("Device         : %s\n", result.Product.Name)
	} else if result.Model != nil {
		fmt.Printf("Device         : %s\n", result.Model.Name)
	} else if result.Valid {
		fmt.Println("Device         : Unknown")
	}

	if result.Model != nil {
		fmt.Printf("Model Identifier: %s\n", result.Model.Model)

		if result.Model.ModelYear != 0 {
			fmt.Printf("Model Year     : %d\n", result.Model.ModelYear)
		}
	}

	fmt.Printf("Reason         : %s\n", result.Reason)
	fmt.Println("================================")
}
