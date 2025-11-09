package main

import (
	"fmt"
	"go-armory/fp/monad"
)

func example2_Configuration() {
	fmt.Println("=== Example 2: Configuration ===")

	config := map[string]string{
		"host": "localhost",
		"port": "8080",
		// "timeout" is missing
	}

	host := monad.GetFromMap(config, "host").GetOrElse("0.0.0.0")
	port := monad.GetFromMap(config, "port").GetOrElse("3000")
	timeout := monad.GetFromMap(config, "timeout").GetOrElse("30")

	fmt.Printf("Server Config:\n")
	fmt.Printf("  Host: %s\n", host)
	fmt.Printf("  Port: %s\n", port)
	fmt.Printf("  Timeout: %s (default)\n", timeout)

	fmt.Println()
}

func main() {
	fmt.Println("Maybe Monad Practical Examples")
	fmt.Println("================================\n")

	example2_Configuration()
	//	example4_QueryParameters()

	//	example5_FeatureFlags()
	//	example6_Cache()
	//	example7_ParsingChain()
	//	example8_Environment()
	//	example9_NestedConfig()
	//	example10_PointerValues()

	fmt.Println("All examples completed!")
}

// Example 1: User Database Lookup
type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

var userDB = map[int]User{
	1: {1, "Alice", "alice@example.com", 30},
	2: {2, "Bob", "bob@example.com", 25},
	3: {3, "Charlie", "charlie@example.com", 35},
}

// Example 2: Configuration with Defaults
// Example 3: Chaining with fmap (not exposed, but showing the concept)
func toUpperCase(s string) string {
	return fmt.Sprintf("%s", s) // Simplified for example
}

// Example 4: HTTP Query Parameters
func example4_QueryParameters() {
	fmt.Println("=== Example 4: Query Parameters ===")

	// Simulate URL: /api/users?page=2&limit=50
	params := map[string]string{
		"page":  "2",
		"limit": "50",
		// "sort" is missing
	}

	pageStr := monad.GetFromMap(params, "page").GetOrElse("1")
	page := monad.ParseNumber(pageStr).GetOrElse(1)

	limitStr := monad.GetFromMap(params, "limit").GetOrElse("10")
	limit := monad.ParseNumber(limitStr).GetOrElse(10)

	sort := monad.GetFromMap(params, "sort").GetOrElse("name")

	fmt.Printf("Query Parameters:\n")
	fmt.Printf("  Page: %d\n", page)
	fmt.Printf("  Limit: %d\n", limit)
	fmt.Printf("  Sort: %s (default)\n", sort)

	// Test with invalid number
	badParams := map[string]string{
		"page": "abc", // invalid number
	}
	badPageStr := monad.GetFromMap(badParams, "page").GetOrElse("1")
	badPage := monad.ParseNumber(badPageStr).GetOrElse(1)
	fmt.Printf("  Invalid page 'abc' -> %d (default)\n", badPage)

	fmt.Println()
}

// Example 5: Feature Flags
func example5_FeatureFlags() {
	fmt.Println("=== Example 5: Feature Flags ===")

	features := map[string]bool{
		"dark_mode":     true,
		"beta_features": false,
		"new_ui":        true,
	}

	darkMode := monad.GetFromMap(features, "dark_mode").GetOrElse(false)
	betaFeatures := monad.GetFromMap(features, "beta_features").GetOrElse(false)
	analytics := monad.GetFromMap(features, "analytics").GetOrElse(true)

	fmt.Printf("Feature Flags:\n")
	fmt.Printf("  Dark Mode: %v\n", darkMode)
	fmt.Printf("  Beta Features: %v\n", betaFeatures)
	fmt.Printf("  Analytics: %v (default, not configured)\n", analytics)

	fmt.Println()
}

// Example 6: Cache with GetOrElse pattern
type Cache struct {
	data map[string]string
}

func (c *Cache) Get(key string) monad.Maybe[string] {
	return monad.GetFromMap(c.data, key)
}

func (c *Cache) GetOrCompute(key string, compute func() string) string {
	cached := c.Get(key)

	switch cached.(type) {
	case monad.JustMaybe[string]:
		fmt.Printf("  [CACHE HIT] %s\n", key)
		return cached.Get()
	case monad.NothingMaybe[string]:
		fmt.Printf("  [CACHE MISS] %s - computing...\n", key)
		value := compute()
		c.data[key] = value
		return value
	default:
		return compute()
	}
}

func example6_Cache() {
	fmt.Println("=== Example 6: Cache ===")

	cache := &Cache{
		data: map[string]string{
			"user:1": "Alice",
		},
	}

	// Hit cache
	name1 := cache.GetOrCompute("user:1", func() string {
		return "Expensive DB Query for user:1"
	})
	fmt.Printf("Result: %s\n", name1)

	// Miss cache
	name2 := cache.GetOrCompute("user:2", func() string {
		return "Bob"
	})
	fmt.Printf("Result: %s\n", name2)

	// Hit cache again (now cached from previous miss)
	name2Again := cache.GetOrCompute("user:2", func() string {
		return "This won't be called"
	})
	fmt.Printf("Result: %s\n", name2Again)

	fmt.Println()
}

// Example 7: Parsing with Fallback Chain
func example7_ParsingChain() {
	fmt.Println("=== Example 7: Parsing with Fallbacks ===")

	inputs := []string{"42", "invalid", "100", ""}

	for _, input := range inputs {
		numMaybe := monad.ParseNumber(input)

		switch numMaybe.(type) {
		case monad.JustMaybe[int]:
			fmt.Printf("'%s' -> %d ✓\n", input, numMaybe.Get())
		case monad.NothingMaybe[int]:
			defaultVal := numMaybe.GetOrElse(0)
			fmt.Printf("'%s' -> %d (default) ✗\n", input, defaultVal)
		}
	}

	fmt.Println()
}

// Example 8: Environment-like Configuration
func example8_Environment() {
	fmt.Println("=== Example 8: Environment Variables ===")

	env := map[string]string{
		"APP_NAME":    "MyApp",
		"PORT":        "8080",
		"DEBUG":       "true",
		"MAX_WORKERS": "10",
		// DATABASE_URL is missing
	}

	appName := monad.GetFromMap(env, "APP_NAME").GetOrElse("DefaultApp")
	port := monad.GetFromMap(env, "PORT").GetOrElse("3000")
	dbURL := monad.GetFromMap(env, "DATABASE_URL").GetOrElse("localhost:5432")

	maxWorkersStr := monad.GetFromMap(env, "MAX_WORKERS").GetOrElse("5")
	maxWorkers := monad.ParseNumber(maxWorkersStr).GetOrElse(5)

	fmt.Printf("Environment:\n")
	fmt.Printf("  APP_NAME: %s\n", appName)
	fmt.Printf("  PORT: %s\n", port)
	fmt.Printf("  DATABASE_URL: %s (default)\n", dbURL)
	fmt.Printf("  MAX_WORKERS: %d\n", maxWorkers)

	fmt.Println()
}

// Example 9: Multi-level Map Access
func example9_NestedConfig() {
	fmt.Println("=== Example 9: Nested Configuration ===")

	// Simulate nested config structure
	databaseConfig := map[string]string{
		"host":     "db.example.com",
		"port":     "5432",
		"database": "myapp",
	}

	serverConfig := map[string]string{
		"host": "api.example.com",
		"port": "8080",
	}

	configs := map[string]map[string]string{
		"database": databaseConfig,
		"server":   serverConfig,
	}

	// Access nested config safely
	dbConfigMaybe := monad.GetFromMap(configs, "database")

	switch dbConfigMaybe.(type) {
	case monad.JustMaybe[map[string]string]:
		dbConfig := dbConfigMaybe.Get()
		host := monad.GetFromMap(dbConfig, "host").GetOrElse("localhost")
		port := monad.GetFromMap(dbConfig, "port").GetOrElse("5432")
		fmt.Printf("Database: %s:%s\n", host, port)
	case monad.NothingMaybe[map[string]string]:
		fmt.Println("Database config not found")
	}

	// Try to access non-existent config
	cacheConfigMaybe := monad.GetFromMap(configs, "cache")
	switch cacheConfigMaybe.(type) {
	case monad.JustMaybe[map[string]string]:
		fmt.Println("Cache config found")
	case monad.NothingMaybe[map[string]string]:
		fmt.Println("Cache config not found - using defaults")
	}

	fmt.Println()
}

// Example 10: Type-safe Pointer Handling
func example10_PointerValues() {
	fmt.Println("=== Example 10: Pointer Values ===")

	val1 := 42
	val2 := 100

	pointerMap := map[string]*int{
		"answer": &val1,
		"max":    &val2,
		"null":   nil, // nil pointer but key exists!
	}

	// Get existing pointer
	answerMaybe := monad.GetFromMap(pointerMap, "answer")
	switch answerMaybe.(type) {
	case monad.JustMaybe[*int]:
		ptr := answerMaybe.Get()
		if ptr != nil {
			fmt.Printf("answer: %d\n", *ptr)
		} else {
			fmt.Println("answer: key exists but value is nil")
		}
	case monad.NothingMaybe[*int]:
		fmt.Println("answer: key not found")
	}

	nullMaybe := monad.GetFromMap(pointerMap, "null")

	switch nullMaybe.(type) {
	case monad.JustMaybe[*int]:
		ptr := nullMaybe.Get()
		if ptr == nil {
			fmt.Println("null: key exists, value is nil pointer")
		}
	case monad.NothingMaybe[*int]:
		fmt.Println("null: key not found")
	}
	missingMaybe := monad.GetFromMap(pointerMap, "missing")

	switch missingMaybe.(type) {
	case monad.JustMaybe[*int]:
		fmt.Println("missing: found")
	case monad.NothingMaybe[*int]:
		fmt.Println("missing: key not found")
	}

	fmt.Println()
}
