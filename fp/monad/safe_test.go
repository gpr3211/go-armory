package monad

import (
	"testing"
)

func TestGetFromMap_KeyExists(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := GetFromMap(m, "a")

	if _, ok := result.(JustMaybe[int]); !ok {
		t.Fatal("Expected JustMaybe")
	}
	if result.Get() != 1 {
		t.Errorf("Expected 1, got %d", result.Get())
	}
}

func TestGetFromMap_KeyDoesNotExist(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := GetFromMap(m, "c")

	if _, ok := result.(NothingMaybe[int]); !ok {
		t.Fatal("Expected NothingMaybe")
	}
}

func TestGetFromMap_EmptyMap(t *testing.T) {
	m := map[string]int{}
	result := GetFromMap(m, "a")

	if _, ok := result.(NothingMaybe[int]); !ok {
		t.Fatal("Expected NothingMaybe for empty map")
	}
}

func TestGetFromMap_NilMap(t *testing.T) {
	var m map[string]int // nil map
	result := GetFromMap(m, "a")

	if _, ok := result.(NothingMaybe[int]); !ok {
		t.Fatal("Expected NothingMaybe for nil map")
	}
}

func TestGetFromMap_ZeroValue(t *testing.T) {
	// Edge case: value is the zero value of its type
	m := map[string]int{"zero": 0}
	result := GetFromMap(m, "zero")

	if _, ok := result.(JustMaybe[int]); !ok {
		t.Fatal("Expected JustMaybe for zero value")
	}
	if result.Get() != 0 {
		t.Errorf("Expected 0, got %d", result.Get())
	}
}

func TestGetFromMap_EmptyString(t *testing.T) {
	// Empty string as value
	m := map[string]string{"empty": ""}
	result := GetFromMap(m, "empty")

	if _, ok := result.(JustMaybe[string]); !ok {
		t.Fatal("Expected JustMaybe for empty string")
	}
	if result.Get() != "" {
		t.Errorf("Expected empty string, got %s", result.Get())
	}
}

func TestGetFromMap_PointerValue(t *testing.T) {
	val := 42
	m := map[string]*int{"ptr": &val}
	result := GetFromMap(m, "ptr")

	if _, ok := result.(JustMaybe[*int]); !ok {
		t.Fatal("Expected JustMaybe")
	}
	if *result.Get() != 42 {
		t.Errorf("Expected 42, got %d", *result.Get())
	}
}

func TestGetFromMap_NilPointerValue(t *testing.T) {
	// Nil pointer as a value (exists in map, but is nil)
	m := map[string]*int{"nilptr": nil}
	result := GetFromMap(m, "nilptr")

	if _, ok := result.(JustMaybe[*int]); !ok {
		t.Fatal("Expected JustMaybe (key exists even though value is nil)")
	}
	if result.Get() != nil {
		t.Error("Expected nil pointer")
	}
}

func TestGetFromMap_StructValue(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	m := map[int]User{
		1: {"Alice", 30},
		2: {"Bob", 25},
	}
	result := GetFromMap(m, 1)

	if _, ok := result.(JustMaybe[User]); !ok {
		t.Fatal("Expected JustMaybe")
	}
	user := result.Get()
	if user.Name != "Alice" || user.Age != 30 {
		t.Errorf("Expected Alice, 30, got %s, %d", user.Name, user.Age)
	}
}

func TestGetFromMap_GetOrElse(t *testing.T) {
	m := map[string]int{"a": 1}

	// Key exists
	result1 := GetFromMap(m, "a")
	if result1.GetOrElse(99) != 1 {
		t.Error("GetOrElse should return actual value when key exists")
	}

	// Key doesn't exist
	result2 := GetFromMap(m, "b")
	if result2.GetOrElse(99) != 99 {
		t.Error("GetOrElse should return default value when key doesn't exist")
	}
}

func TestGetFromMap_ChainWithFmap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	// Test fmap chaining with existing key
	result1 := fmap(GetFromMap(m, "a"), func(x int) int { return x * 2 })
	if result1.Get() != 2 {
		t.Errorf("Expected 2, got %d", result1.Get())
	}

	// Test fmap chaining with non-existing key
	result2 := fmap(GetFromMap(m, "c"), func(x int) int { return x * 2 })
	if _, ok := result2.(NothingMaybe[int]); !ok {
		t.Fatal("Expected NothingMaybe after fmap on Nothing")
	}
}

func TestGetFromMap_InterfaceValue(t *testing.T) {
	m := map[string]any{"key": "string value", "num": 42}

	result := GetFromMap(m, "key")
	if _, ok := result.(JustMaybe[any]); !ok {
		t.Fatal("Expected JustMaybe")
	}
	if result.Get().(string) != "string value" {
		t.Error("Unexpected interface value")
	}
}

func TestGetFromMap_BooleanKey(t *testing.T) {
	m := map[bool]string{true: "yes", false: "no"}

	result := GetFromMap(m, true)
	if result.Get() != "yes" {
		t.Errorf("Expected 'yes', got %s", result.Get())
	}
}
