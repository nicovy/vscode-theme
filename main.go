package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

// brightenColor takes a hex color string and increases its brightness by the given percentage.
func brightenColor(hex string, percent float64) (string, error) {
	var r, g, b int
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return "", fmt.Errorf("invalid hex color format")
	}

	// Increase brightness by the percentage
	r = int(math.Min(float64(r)*(1+percent/100), 255))
	g = int(math.Min(float64(g)*(1+percent/100), 255))
	b = int(math.Min(float64(b)*(1+percent/100), 255))

	// Return the new hex color
	return fmt.Sprintf("#%02x%02x%02x", r, g, b), nil
}

// brightnessDifference calculates the percentage difference in brightness between two hex color strings.
// It returns the average percentage increase across RGB channels from baseHex to brighterHex.
// Returns a positive value if brighterHex is brighter, negative if darker, or an error if invalid hex formats.
func brightnessDifference(baseHex string, brighterHex string) (float64, error) {
	var r1, g1, b1, r2, g2, b2 int
	_, err := fmt.Sscanf(baseHex, "#%02x%02x%02x", &r1, &g1, &b1)
	if err != nil {
		return -1, fmt.Errorf("invalid hex color format")
	}
	_, err = fmt.Sscanf(brighterHex, "#%02x%02x%02x", &r2, &g2, &b2)
	if err != nil {
		return -1, fmt.Errorf("invalid hex color format")
	}

	// Calculate the brightness difference
	percentR := math.Round(((float64(r2) / float64(r1)) - 1) * 100)
	percentG := math.Round(((float64(g2) / float64(g1)) - 1) * 100)
	percentB := math.Round(((float64(b2) / float64(b1)) - 1) * 100)
	avg := math.Round((percentR + percentG + percentB) / 3)

	return avg, nil
}

// loads a json file
func loadJson(filename string) (map[string]interface{}, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	// Parse JSON content
	var result map[string]interface{}
	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return result, nil
}

func main() {
	// load themes-color.json
	themes, err := loadJson("themes-color.json")
	if err != nil {
		fmt.Printf("Error loading themes: %v\n", err)
		return
	}

	themesList, ok := themes["themes"].([]interface{})
	if !ok {
		fmt.Println("Error: Could not parse themes array")
		return
	}

	packageJson, err := loadJson("package.json")
	if err != nil {
		fmt.Printf("Error loading package.json: %v\n", err)
		return
	}
	shortName, ok := packageJson["shortName"].(string)
	if !ok {
		fmt.Println("Error: Could not parse shortName from package.json")
		return
	}
	themesArray := make([]map[string]interface{}, 0)

	// load template/editor-dark-color-theme.json
	templateSyntax, err := loadJson("template/syntax-dark-color-theme.json")
	if err != nil {
		fmt.Printf("Error loading syntax template: %v\n", err)
		return
	}

	// load template/editor-dark-color-theme.json
	templateEditor, err := loadJson("template/editor-dark-color-theme.json")
	if err != nil {
		fmt.Printf("Error loading editor template: %v\n", err)
		return
	}
	// combine both objects
	combined := make(map[string]interface{})
	for k, v := range templateSyntax {
		combined[k] = v
	}
	for k, v := range templateEditor {
		combined[k] = v
	}
	// convert combined to json string
	combinedJsonDark, err := json.Marshal(combined)
	if err != nil {
		fmt.Printf("Error marshalling combined JSON: %v\n", err)
		return
	}

	// TODO!: load template/editor-light-color-theme.json
	templateSyntax, err = loadJson("template/syntax-dark-color-theme.json")
	if err != nil {
		fmt.Printf("Error loading syntax template: %v\n", err)
		return
	}

	// TODO!: load template/editor-light-color-theme.json
	templateEditor, err = loadJson("template/editor-dark-color-theme.json")
	if err != nil {
		fmt.Printf("Error loading editor template: %v\n", err)
		return
	}
	// combine both objects
	combined = make(map[string]interface{})
	for k, v := range templateSyntax {
		combined[k] = v
	}
	for k, v := range templateEditor {
		combined[k] = v
	}
	combinedJsonLight, err := json.Marshal(combined)
	if err != nil {
		fmt.Printf("Error marshalling combined JSON: %v\n", err)
		return
	}

	// Access the themes array

	// Process themes
	for _, theme := range themesList {
		themeMap, ok := theme.(map[string]interface{})
		if !ok {
			continue
		}

		themeType := themeMap["type"].(string)
		name := themeMap["name"].(string)
		baseColor := themeMap["baseColor"].(string)
		accentColor := themeMap["accentColor"].(string)
		sColor := themeMap["secondaryColor"]
		bLevel := themeMap["brighterLevel"]
		iColor := themeMap["inputColor"]

		var inputColor string
		var secondaryColor string
		brighterLevel := -1.0

		if iColor != nil {
			inputColor = iColor.(string)
		} else {
			inputColor, err = brightenColor(baseColor, -15)
		}
		if sColor != nil {
			secondaryColor = sColor.(string)
			brighterLevel, err = brightnessDifference(baseColor, secondaryColor)
		} else {
			brighterLevel = 35
			if themeType == "light" {
				brighterLevel = -10
			}
			if bLevel != nil {
				brighterLevel = bLevel.(float64)
			}
			secondaryColor, err = brightenColor(baseColor, brighterLevel)
		}

		var combinedString string
		if themeType == "dark" {
			combinedString = string(combinedJsonDark)
		} else {
			combinedString = string(combinedJsonLight)
			continue // TODO!: remove this when light theme is added
		}
		combinedString = strings.ReplaceAll(combinedString, "[BASE]", baseColor)
		combinedString = strings.ReplaceAll(combinedString, "[BRIGHT]", secondaryColor)
		combinedString = strings.ReplaceAll(combinedString, "[ACCENT]", accentColor)
		combinedString = strings.ReplaceAll(combinedString, "[INPUT]", inputColor)

		// write to file
		fileName := strings.ReplaceAll(fmt.Sprintf("themes/%s-%s-color-theme.json", strings.ToLower(themeType), strings.ToLower(name)), " ", "-")

		err = os.WriteFile(fileName, []byte(combinedString), 0644)
		if err != nil {
			fmt.Printf("Error writing file %s: %v\n", fileName, err)
			return
		}
		fmt.Printf("Theme: %s, %s, Base: %s, Highlight: %s, Secondary: %s, Brighter Level: %f\n", themeType, name, baseColor, accentColor, secondaryColor, brighterLevel)

		uiTheme := "vs"
		if themeType == "dark" {
			uiTheme = "vs-dark"
		}
		themesArray = append(themesArray, map[string]interface{}{
			"label":   fmt.Sprintf("%s - %s", shortName, strings.ToLower(name)),
			"uiTheme": uiTheme,
			"path":    fmt.Sprintf("./%s", fileName)})
	}

	themesObject := map[string]interface{}{
		"themes": themesArray,
	}
	packageJson["contributes"] = themesObject
	packageToJson, err := json.MarshalIndent(packageJson, "", "  ")
	if err != nil {
		fmt.Printf("Error converting to json package.json: %v\n", err)
		return
	}
	err = os.WriteFile("package.json", []byte(packageToJson), 0644)
	if err != nil {
		fmt.Printf("Error writing file package.json: %v\n", err)
		return
	}
}
