package utls

import (
	"testing"
)

// Original implementation
func PermutationsOriginal(s string) []string {
	var result []string
	chars := []rune(s)
	var backtrack func(pos int)
	backtrack = func(pos int) {
		if pos == len(chars) {
			result = append(result, string(chars))
			return
		}
		for i := pos; i < len(chars); i++ {
			chars[pos], chars[i] = chars[i], chars[pos]
			backtrack(pos + 1)
			chars[pos], chars[i] = chars[i], chars[pos]
		}
	}
	backtrack(0)
	unique := make(map[string]bool, len(result))
	out := []string{}
	for _, v := range result {
		_, ok := unique[v]
		if !ok {
			unique[v] = true
			out = append(out, v)
		}
	}
	return out
}

// Optimized implementation
func PermutationsOptimized(s string) []string {
	chars := []rune(s)
	n := len(chars)

	// Calculate approximate capacity
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

// Benchmarks
func BenchmarkPermutationsOriginal_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PermutationsOriginal("abc")
	}
}

func BenchmarkPermutationsOptimized_Short(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PermutationsOptimized("abc")
	}
}

func BenchmarkPermutationsOriginal_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PermutationsOriginal("abcd12")
	}
}

func BenchmarkPermutationsOptimized_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PermutationsOptimized("abcd12")
	}
}

func BenchmarkPermutationsOriginal_WithDuplicates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PermutationsOriginal("aabbcc")
	}
}

func BenchmarkPermutationsOptimized_WithDuplicates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PermutationsOptimized("aabbcc")
	}
}

func BenchmarkPermutationsOriginal_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PermutationsOriginal("abcd8391")
	}
}

func BenchmarkPermutationsOptimized_Long(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PermutationsOptimized("abcd8391")
	}
}

// Test to verify both produce same results
func TestPermutationsMatch(t *testing.T) {
	testCases := []string{
		"abc",
		"aab",
		"abcd",
		"aabbcc",
		"abcd8391",
	}

	for _, tc := range testCases {
		orig := PermutationsOriginal(tc)
		opt := PermutationsOptimized(tc)

		if len(orig) != len(opt) {
			t.Errorf("Length mismatch for %q: original=%d, optimized=%d", tc, len(orig), len(opt))
		}

		// Convert to maps for comparison (order doesn't matter)
		origMap := make(map[string]bool)
		for _, v := range orig {
			origMap[v] = true
		}

		optMap := make(map[string]bool)
		for _, v := range opt {
			optMap[v] = true
		}

		for k := range origMap {
			if !optMap[k] {
				t.Errorf("Original has %q but optimized doesn't for input %q", k, tc)
			}
		}

		for k := range optMap {
			if !origMap[k] {
				t.Errorf("Optimized has %q but original doesn't for input %q", k, tc)
			}
		}
	}
}
