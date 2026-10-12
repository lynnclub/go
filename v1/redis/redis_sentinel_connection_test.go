package redis

import (
	"strings"
	"testing"
)

func TestSentinelConnectionFailure(t *testing.T) {
	for _, tc := range []struct {
		name string
		tls  bool
	}{
		{name: "without_tls"},
		{name: "with_tls", tls: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := "test_sentinel_connection_" + tc.name
			options[name] = Option{
				Address:    []string{"127.0.0.1:0"},
				MasterName: "mymaster",
				TLS:        tc.tls,
			}
			t.Cleanup(func() {
				delete(options, name)
				poolSentinel.Delete(name)
			})

			defer func() {
				got, ok := recover().(string)
				if !ok || !strings.Contains(got, "Failed to connect redis sentinel "+name+" err:") {
					t.Fatalf("expected connection error panic, got %v", got)
				}
			}()

			Sentinel(name)
		})
	}
}
