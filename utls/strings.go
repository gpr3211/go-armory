package utls

import "unicode"

// IsUpperCase returns true only if entire string is in uppercase.
func IsUpperCase(s string) bool {
	letter := false
	for _, v := range s {
		if unicode.IsLetter(v) {
			letter = true
			if unicode.IsLower(v) {
				return false
			}
		}
	}
	return letter
}
