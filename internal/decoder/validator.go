package decoder

import "strconv"

const base34 = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func isBase34(c byte) bool {
	for i := 0; i < len(base34); i++ {
		if base34[i] == c {
			return true
		}
	}

	return false
}

func validateBase34(value string) bool {
	for i := 0; i < len(value); i++ {
		if !isBase34(value[i]) {
			return false
		}
	}

	return true
}

func decode11(serial string, result Result) Result {
	result.Format = FormatLegacy11

	result.Location = serial[0:2]

	if location, found := lookupLocation(
		result.Location,
		"legacy-11",
	); found {
		result.LocationInfo = &location
	}

	result.YearCode = serial[2:3]
	result.WeekCode = serial[3:5]
	result.Production = serial[5:8]
	result.ProductCode = serial[8:11]

	if serial[2] < '0' || serial[2] > '9' {
		result.Reason = "Invalid legacy year code."
		return result
	}

	result.Year = 2000 + int(serial[2]-'0')

	week, err := strconv.Atoi(serial[3:5])
	if err != nil || week < 1 || week > 53 {
		result.Reason = "Invalid production week."
		return result
	}

	result.Week = week

	if !validateBase34(result.Production) {
		result.Reason = "Invalid production identifier."
		return result
	}

	if !validateBase34(result.ProductCode) {
		result.Reason = "Invalid product/configuration code."
		return result
	}

	if product, found := lookupProduct(result.ProductCode); found {
		result.Product = &product
	}

	if model, found := lookupModel(result.ProductCode); found {
		result.Model = &model
	}

	result.Valid = true

	if result.Product != nil || result.Model != nil {
		result.Reason = "Valid Apple legacy 11-character structure. Product recognized."
	} else {
		result.Reason = "Valid Apple legacy 11-character structure. Product code not found."
	}

	return result
}

func decode12(serial string, result Result) Result {
	result.Format = FormatLegacy12

	result.Location = serial[0:3]

	if location, found := lookupLocation(
		result.Location,
		"legacy-12",
	); found {
		result.LocationInfo = &location
	}

	result.YearCode = serial[3:4]
	result.WeekCode = serial[4:5]
	result.Production = serial[5:8]
	result.ProductCode = serial[8:12]

	// Resolve product and model before resolving the year.
	if product, found := lookupProduct(result.ProductCode); found {
		result.Product = &product
	}

	if model, found := lookupModel(result.ProductCode); found {
		result.Model = &model
	}

	year, half, ok := decodeYear(
		serial[3],
		result.Model,
	)

	if !ok {
		result.Reason = "Invalid legacy year code."
		return result
	}

	result.Year = year
	result.Half = half

	week, ok := decodeWeek(
		serial[4],
		serial[3],
	)

	if !ok {
		result.Reason = "Invalid production week code."
		return result
	}

	result.Week = week

	if !validateBase34(result.Production) {
		result.Reason = "Invalid production identifier."
		return result
	}

	if !validateBase34(result.ProductCode) {
		result.Reason = "Invalid product/configuration code."
		return result
	}

	result.Valid = true

	if result.Product != nil || result.Model != nil {
		result.Reason = "Valid Apple legacy 12-character structure. Product recognized."
	} else {
		result.Reason = "Valid Apple legacy 12-character structure. Product code not found."
	}

	return result
}

func decodeYear(c byte, model *ModelInfo) (int, int, bool) {
	yearTable := map[byte][2]int{
		'C': {2010, 2020},
		'D': {2010, 2020},
		'F': {2011, 2021},
		'G': {2011, 2021},
		'H': {2012, 2022},
		'J': {2012, 2022},
		'K': {2013, 2023},
		'L': {2013, 2023},
		'M': {2014, 2024},
		'N': {2014, 2024},
		'P': {2015, 2025},
		'Q': {2015, 2025},
		'R': {2016, 2026},
		'S': {2016, 2026},
		'T': {2017, 2027},
		'V': {2017, 2027},
		'W': {2018, 2028},
		'X': {2018, 2028},
		'Y': {2019, 2029},
		'Z': {2019, 2029},
	}

	pair, ok := yearTable[c]
	if !ok {
		return 0, 0, false
	}

	half := 1

	switch c {
	case 'D', 'G', 'J', 'L', 'N', 'Q', 'S', 'V', 'X', 'Z':
		half = 2
	}

	// If the identified model belongs to the 2020s,
	// use the second year represented by the serial code.
	if model != nil && model.ModelYear >= 2020 {
		return pair[1], half, true
	}

	// Legacy devices remain in the original decade.
	return pair[0], half, true
}

func decodeWeek(weekCode byte, yearCode byte) (int, bool) {
	weekTable := map[byte]int{
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'C': 10,
		'D': 11,
		'F': 12,
		'G': 13,
		'H': 14,
		'J': 15,
		'K': 16,
		'L': 17,
		'M': 18,
		'N': 19,
		'P': 20,
		'Q': 21,
		'R': 22,
		'T': 23,
		'V': 24,
		'W': 25,
		'X': 26,
		'Y': 27,
	}

	week, ok := weekTable[weekCode]
	if !ok {
		return 0, false
	}

	switch yearCode {
	case 'D', 'G', 'J', 'L', 'N', 'Q', 'S', 'V', 'X', 'Z':
		week += 26
	}

	if week < 1 || week > 53 {
		return 0, false
	}

	return week, true
}
