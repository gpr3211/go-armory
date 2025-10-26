package utls

import "unicode"

// IsUpperCase returns true only if the entire string is in uppercase.
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
func Permutations(s string) []string {
	chars := []rune(s)
	n := len(chars)

	// Calculate approximate capacity (won't be exact with duplicates)
	capacity := 1
	for i := 2; i <= n && i <= 12; i++ {
		capacity *= i
	}

	unique := make(map[string]bool, capacity)

	var backtrack func(pos int)
	backtrack = func(pos int) {
		if pos == n {
			unique[string(chars)] = true
			return
		}

		// Use a set to avoid duplicate swaps at this level
		seen := make(map[rune]bool)
		for i := pos; i < n; i++ {
			// Skip if we've already used this character at this position
			if seen[chars[i]] {
				continue
			}
			seen[chars[i]] = true

			chars[pos], chars[i] = chars[i], chars[pos]
			backtrack(pos + 1)
			chars[pos], chars[i] = chars[i], chars[pos]
		}
	}

	backtrack(0)

	// Convert map to slice with pre-allocated capacity
	out := make([]string, 0, len(unique))
	for v := range unique {
		out = append(out, v)
	}
	return out
}
