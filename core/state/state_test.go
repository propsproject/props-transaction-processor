package state

import "testing"

func TestNewState(t *testing.T) {
	if s := NewState(nil); s == nil {
		t.Fatalf("expected state instance to be defined for (%v)", s)
	}
}
