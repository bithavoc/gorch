package gorch

import (
	"testing"
)

func TestOperationsRegistry(t *testing.T) {
	registry := &OperationsRegistry{}
	entry := registry.Register("sum")
	if !registry.Exists("sum") {
		t.Fatalf("Registry should report sum as existent")
	}
	if entry.IsHosted() {
		t.Fatalf("Operation entry should not be reported as Hosted yet")
	}
	entry.Host(func(inv Invocation) (interface{}, error) { return nil, nil })
	if !entry.IsHosted() {
		t.Fatalf("Operation entry should be reported as hosted")
	}
}
