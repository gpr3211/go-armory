package utls

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Takes JSON and prints it out with indents and colored keys
func PrettyPrint(v any) (err error) {
	// Marshal with indentation
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	// Convert JSON to string and use a regex to match keys
	jsonStr := string(b)
	re := regexp.MustCompile(`"([^"]+)"\s*:`)

	// Replace each key with a blue-colored version
	coloredStr := re.ReplaceAllStringFunc(jsonStr, func(match string) string {
		// Color the key in blue and return it
		return fmt.Sprintf("\033[34m%s\033[0m", match)
	})

	// Print the colored JSON
	fmt.Println(coloredStr)
	return nil
}
