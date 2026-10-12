package redis

import (
	"strings"
	"testing"
)

func TestClusterConnectionFailure(t *testing.T) {
	for _, tc := range []struct {
		name string
		tls  bool
	}{
		{name: "without_tls"},
		{name: "with_tls", tls: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := "test_cluster_connection_" + tc.name
			options[name] = Option{
				Address: []string{"127.0.0.1:0"},
				TLS:     tc.tls,
			}
			t.Cleanup(func() {
				delete(options, name)
				poolCluster.Delete(name)
			})

			defer func() {
				got, ok := recover().(string)
				if !ok || !strings.Contains(got, "Failed to connect redis cluster "+name+" err:") {
					t.Fatalf("expected connection error panic, got %v", got)
				}
			}()

			Cluster(name)
		})
	}
}
