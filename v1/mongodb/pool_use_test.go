package mongodb

import (
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestUseCachedClient(t *testing.T) {
	client := &mongo.Client{}
	pool.Store("default", client)
	t.Cleanup(func() {
		pool.Delete("default")
	})

	if got := Use(""); got != client {
		t.Fatal("Use did not reuse the default client")
	}
}

func TestUseMissingOptionPanics(t *testing.T) {
	const name = "test_use_missing_option"
	delete(options, name)
	t.Cleanup(func() {
		delete(options, name)
	})

	defer func() {
		if got := recover(); got != "Option not found "+name {
			t.Fatalf("expected missing-option panic, got %v", got)
		}
	}()

	Use(name)
}

func TestUseInvalidURIIsReported(t *testing.T) {
	const name = "test_use_invalid_uri"
	options[name] = Option{DSN: "invalid://"}
	t.Cleanup(func() {
		delete(options, name)
		pool.Delete(name)
	})

	defer func() {
		got, ok := recover().(string)
		if !ok || !strings.Contains(got, "Failed to connect mongodb "+name+" err:") {
			t.Fatalf("expected connection error panic, got %v", got)
		}
	}()

	Use(name)
}
