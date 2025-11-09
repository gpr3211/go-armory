package main

import (
	"fmt"
	"go-armory/fp/monad"
	"math/rand"
	"time"
)

// Example 1: Basic Future Creation and Getting Results
func example1_BasicFuture() {
	fmt.Println("=== Example 1: Basic Future ===")

	future := monad.NewFuture(func() (int, error) {
		fmt.Println("  Computing in background...")
		time.Sleep(1 * time.Second)
		return 42, nil
	})

	fmt.Println("  Future created, continuing main thread...")
	fmt.Println("  Doing other work...")

	result, err := future.Get()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Result: %d\n", result)
	}

	fmt.Println()
}

// Example 2: Future with Timeout
func example2_FutureWithTimeout() {
	fmt.Println("=== Example 2: Future with Timeout ===")

	// executed future
	//
	//
	fastFuture := monad.NewFuture(func() (string, error) {
		time.Sleep(500 * time.Millisecond)
		return "Fast result", nil
	})

	result1, err1, timedOut1 := fastFuture.GetWithTimeout(2 * time.Second)
	if timedOut1 {
		fmt.Println("  Fast future: TIMEOUT")
	} else if err1 != nil {
		fmt.Printf("  Fast future error: %v\n", err1)
	} else {
		fmt.Printf("  Fast future: %s ✓\n", result1)
	}

	// timed out future
	//
	//
	slowFuture := monad.NewFuture(func() (string, error) {
		time.Sleep(3 * time.Second)
		return "This won't be seen", nil
	})

	result2, err2, timedOut2 := slowFuture.GetWithTimeout(1 * time.Second)
	if timedOut2 {
		fmt.Println("  Slow future: TIMEOUT ✗")
	} else if err2 != nil {
		fmt.Printf("  Slow future error: %v\n", err2)
	} else {
		fmt.Printf("  Slow future: %s\n", result2)
	}

	fmt.Println()
}

// Example 3: Mapping Over Futures
func example3_MapTransformation() {
	fmt.Println("=== Example 3: Map Transformation ===")

	future := monad.NewFuture(func() (int, error) {
		fmt.Println("  Computing base value...")
		time.Sleep(500 * time.Millisecond)
		return 10, nil
	})

	// Transform: multiply by 2
	doubled := monad.Map(future, func(x int) int {
		fmt.Println("  Doubling...")
		return x * 2
	})

	// Transform again: convert to string
	stringified := monad.Map(doubled, func(x int) string {
		fmt.Println("  Converting to string...")
		return fmt.Sprintf("The answer is %d", x)
	})

	result, err := stringified.Get()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Final result: %s\n", result)
	}

	fmt.Println()
}

// Example 4: FlatMap for Chaining Dependent Futures
func example4_FlatMapChaining() {
	fmt.Println("=== Example 4: FlatMap Chaining ===")

	// First future: fetch user ID
	userIDFuture := monad.NewFuture(func() (int, error) {
		fmt.Println("  Fetching user ID...")
		time.Sleep(500 * time.Millisecond)
		return 123, nil
	})

	// Second future: use user ID to fetch user name
	userNameFuture := monad.FlatMap(userIDFuture, func(userID int) *monad.Future[string] {
		return monad.NewFuture(func() (string, error) {
			fmt.Printf("  Fetching user name for ID %d...\n", userID)
			time.Sleep(500 * time.Millisecond)
			return "Alice", nil
		})
	})

	// Third future: use name to fetch email
	emailFuture := monad.FlatMap(userNameFuture, func(name string) *monad.Future[string] {
		return monad.NewFuture(func() (string, error) {
			fmt.Printf("  Fetching email for %s...\n", name)
			time.Sleep(500 * time.Millisecond)
			return fmt.Sprintf("%s@example.com", name), nil
		})
	})

	email, err := emailFuture.Get()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Final email: %s\n", email)
	}

	fmt.Println()
}

// Example 5: Error Handling
func example5_ErrorHandling() {
	fmt.Println("=== Example 5: Error Handling ===")

	successFuture := monad.NewFuture(func() (int, error) {
		time.Sleep(300 * time.Millisecond)
		return 100, nil
	})

	errorFuture := monad.NewFuture(func() (int, error) {
		time.Sleep(300 * time.Millisecond)
		return 0, fmt.Errorf("database connection failed")
	})

	result1, err1 := successFuture.Get()
	if err1 != nil {
		fmt.Printf("  Success future error: %v\n", err1)
	} else {
		fmt.Printf("  Success future: %d ✓\n", result1)
	}

	result2, err2 := errorFuture.Get()
	if err2 != nil {
		fmt.Printf("  Error future: %v ✗\n", err2)
	} else {
		fmt.Printf("  Error future: %d\n", result2)
	}

	fmt.Println()
}

// Example 6: Immediate Futures (Successful/Failed)
func example6_ImmediateFutures() {
	fmt.Println("=== Example 6: Immediate Futures ===")

	// Already completed future
	success := monad.Successful(42)
	result, err := success.Get()
	fmt.Printf("  Successful: %d (err: %v)\n", result, err)

	// Already failed future
	failed := monad.Failed[string](fmt.Errorf("something went wrong"))
	result2, err2 := failed.Get()
	fmt.Printf("  Failed: '%s' (err: %v)\n", result2, err2)

	fmt.Println()
}

// Example 7: Sequence - Running Multiple Futures in Parallel
func example7_SequenceParallel() {
	fmt.Println("=== Example 7: Sequence (Parallel Execution) ===")

	start := time.Now()

	future1 := monad.NewFuture(func() (int, error) {
		fmt.Println("  Task 1 starting...")
		time.Sleep(1 * time.Second)
		fmt.Println("  Task 1 done")
		return 10, nil
	})

	future2 := monad.NewFuture(func() (int, error) {
		fmt.Println("  Task 2 starting...")
		time.Sleep(1 * time.Second)
		fmt.Println("  Task 2 done")
		return 20, nil
	})

	future3 := monad.NewFuture(func() (int, error) {
		fmt.Println("  Task 3 starting...")
		time.Sleep(1 * time.Second)
		fmt.Println("  Task 3 done")
		return 30, nil
	})

	// Wait for all futures to complete
	allResults := monad.Sequence(future1, future2, future3)
	results, err := allResults.Get()

	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Results: %v\n", results)
		fmt.Printf("  Total time: ~%.1fs\n", elapsed.Seconds())
	}

	fmt.Println()
}

// Example 8: Simulating API Calls
func simulateAPICall(endpoint string, delay time.Duration) *monad.Future[string] {
	return monad.NewFuture(func() (string, error) {
		fmt.Printf("  API call to %s...\n", endpoint)
		time.Sleep(delay)

		// Simulate occasional failures
		if rand.Float32() < 0.2 {
			return "", fmt.Errorf("API error from %s", endpoint)
		}

		return fmt.Sprintf("Data from %s", endpoint), nil
	})
}

func example8_APICallsWithFallback() {
	fmt.Println("=== Example 8: API Calls with Fallback ===")

	rand.Seed(time.Now().UnixNano())

	primaryAPI := simulateAPICall("/api/v1/data", 800*time.Millisecond)

	result, err := primaryAPI.Get()
	if err != nil {
		fmt.Printf("  Primary API failed: %v\n", err)
		fmt.Println("  Trying fallback API...")

		fallbackAPI := simulateAPICall("/api/v2/data", 500*time.Millisecond)
		result, err = fallbackAPI.Get()

		if err != nil {
			fmt.Printf("  Fallback also failed: %v\n", err)
		} else {
			fmt.Printf("  Fallback success: %s\n", result)
		}
	} else {
		fmt.Printf("  Primary success: %s\n", result)
	}

	fmt.Println()
}

// Example 9: Database Query Simulation
type User struct {
	ID    int
	Name  string
	Email string
}

func fetchUserFromDB(id int) *monad.Future[User] {
	return monad.NewFuture(func() (User, error) {
		fmt.Printf("  Querying database for user %d...\n", id)
		time.Sleep(700 * time.Millisecond)

		if id == 999 {
			return User{}, fmt.Errorf("user %d not found", id)
		}

		return User{
			ID:    id,
			Name:  fmt.Sprintf("User%d", id),
			Email: fmt.Sprintf("user%d@example.com", id),
		}, nil
	})
}

func example9_DatabaseQuery() {
	fmt.Println("=== Example 9: Database Query ===")

	userFuture := fetchUserFromDB(42)

	// Transform user to greeting message
	greetingFuture := monad.Map(userFuture, func(u User) string {
		return fmt.Sprintf("Hello, %s! (ID: %d)", u.Name, u.ID)
	})

	greeting, err := greetingFuture.Get()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  %s\n", greeting)
	}

	// Try non-existent user
	fmt.Println("\n  Trying non-existent user...")
	missingUserFuture := fetchUserFromDB(999)
	_, err2 := missingUserFuture.Get()
	if err2 != nil {
		fmt.Printf("  Error: %v ✗\n", err2)
	}

	fmt.Println()
}

// Example 10: Batch Processing Multiple Items
func processItem(id int) *monad.Future[string] {
	return monad.NewFuture(func() (string, error) {
		delay := time.Duration(rand.Intn(500)+300) * time.Millisecond
		time.Sleep(delay)
		return fmt.Sprintf("Processed item %d", id), nil
	})
}

func example10_BatchProcessing() {
	fmt.Println("=== Example 10: Batch Processing ===")

	start := time.Now()

	// Create futures for processing items 1-5
	futures := make([]*monad.Future[string], 5)
	for i := range 5 {
		futures[i] = processItem(i + 1)
		fmt.Printf("  Launched processing for item %d\n", i+1)
	}

	// Wait for all to complete
	allResults := monad.Sequence(futures...)
	results, err := allResults.Get()

	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Println("\n  All items processed:")
		for _, result := range results {
			fmt.Printf("    - %s\n", result)
		}
		fmt.Printf("  Total time: %.2fs (parallel)\n", elapsed.Seconds())
	}

	fmt.Println()
}

// Example 11: Combining Map and FlatMap
func example11_CombiningOperations() {
	fmt.Println("=== Example 11: Combining Map and FlatMap ===")

	// Start with a number
	numberFuture := monad.NewFuture(func() (int, error) {
		time.Sleep(300 * time.Millisecond)
		return 5, nil
	})

	// Double it
	doubledFuture := monad.Map(numberFuture, func(x int) int {
		fmt.Printf("  Doubling %d...\n", x)
		return x * 2
	})

	// Use the doubled value to fetch data
	dataFuture := monad.FlatMap(doubledFuture, func(x int) *monad.Future[string] {
		return monad.NewFuture(func() (string, error) {
			fmt.Printf("  Fetching %d items...\n", x)
			time.Sleep(300 * time.Millisecond)
			return fmt.Sprintf("Fetched %d items", x), nil
		})
	})

	result, err := dataFuture.Get()
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Final result: %s\n", result)
	}

	fmt.Println()
}

// Example 12: Race Condition - First to Complete
func example12_RaceCondition() {
	fmt.Println("=== Example 12: Race Condition (Manual) ===")

	future1 := monad.NewFuture(func() (string, error) {
		time.Sleep(1 * time.Second)
		return "Slow server", nil
	})

	future2 := monad.NewFuture(func() (string, error) {
		time.Sleep(300 * time.Millisecond)
		return "Fast server", nil
	})

	// Manually race using channels
	resultChan := make(chan string, 2)

	go func() {
		result, err := future1.Get()
		if err == nil {
			resultChan <- result
		}
	}()

	go func() {
		result, err := future2.Get()
		if err == nil {
			resultChan <- result
		}
	}()

	winner := <-resultChan
	fmt.Printf("  Winner: %s\n", winner)

	fmt.Println()
}

func main() {
	fmt.Println("Future Monad Practical Examples")
	fmt.Println("=================================")

	example1_BasicFuture()
	example2_FutureWithTimeout()
	example3_MapTransformation()
	example4_FlatMapChaining()
	example5_ErrorHandling()
	example6_ImmediateFutures()
	example7_SequenceParallel()
	example8_APICallsWithFallback()
	example9_DatabaseQuery()
	example10_BatchProcessing()
	example11_CombiningOperations()
	example12_RaceCondition()

	fmt.Println("All examples completed!")
}
