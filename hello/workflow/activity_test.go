package workflow

import (
	"context"
	"testing"
)

func TestComposeGreeting(t *testing.T) {
	ctx := context.Background()
	name := "Alice"
	expected := "Hello Alice!"

	result, err := ComposeGreeting(ctx, name)

	if err != nil {
		t.Fatalf("ComposeGreeting returned an unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected greeting '%s', but got '%s'", expected, result)
	}
}
