package rabbit

import (
	"strings"
	"testing"

	"github.com/wagslane/go-rabbitmq"
)

func TestUseCachedConnection(t *testing.T) {
	conn := &rabbitmq.Conn{}
	pool.Store("test_use_cached", conn)
	t.Cleanup(func() {
		pool.Delete("test_use_cached")
	})

	if got := Use("test_use_cached"); got != conn {
		t.Fatal("Use did not reuse the cached connection")
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

func TestUseInvalidDSNIsReported(t *testing.T) {
	const name = "test_use_invalid_dsn"
	options[name] = Option{DSN: "://"}
	t.Cleanup(func() {
		delete(options, name)
		pool.Delete(name)
	})

	defer func() {
		got, ok := recover().(string)
		if !ok || !strings.Contains(got, "Failed to connect rabbitmq "+name+" err:") {
			t.Fatalf("expected connection error panic, got %v", got)
		}
	}()

	Use(name)
}

func TestGetPublisherReturnsCachedPublisher(t *testing.T) {
	publisher := &rabbitmq.Publisher{}
	poolPublisher.Store("test_publisher_cached", publisher)
	t.Cleanup(func() {
		poolPublisher.Delete("test_publisher_cached")
	})

	if got := GetPublisher("test_publisher_cached"); got != publisher {
		t.Fatal("GetPublisher did not reuse the cached publisher")
	}
}

func TestGetConsumerReturnsCachedConsumer(t *testing.T) {
	consumer := &rabbitmq.Consumer{}
	poolConsumer.Store("test_consumer_cached", consumer)
	t.Cleanup(func() {
		poolConsumer.Delete("test_consumer_cached")
	})

	if got := GetConsumer("test_consumer_cached", "unused"); got != consumer {
		t.Fatal("GetConsumer did not reuse the cached consumer")
	}
}
