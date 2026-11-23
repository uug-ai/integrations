package integrations

import (
	"strings"
	"testing"
)

func TestNewMQTTClientID(t *testing.T) {
	first, err := newMQTTClientID()
	if err != nil {
		t.Fatalf("newMQTTClientID() error = %v", err)
	}
	second, err := newMQTTClientID()
	if err != nil {
		t.Fatalf("newMQTTClientID() error = %v", err)
	}

	if !strings.HasPrefix(first, "hub-pipe-notif-") {
		t.Fatalf("newMQTTClientID() = %q, want hub-pipe-notif- prefix", first)
	}
	if len(first) != 23 {
		t.Fatalf("len(newMQTTClientID()) = %d, want 23", len(first))
	}
	if first == second {
		t.Fatalf("newMQTTClientID() returned duplicate ID %q", first)
	}
}
