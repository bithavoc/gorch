package gorch

import (
	"testing"
)

func TestRegistryAnalysis(t *testing.T) {
	registry := &OperationsRegistry{}
	hosted := registry.Register("hosted")
	hosted.Host(func(inv Invocation) (interface{}, error) { return nil, nil })
	registry.Register("unhosted")
	analysis := AnalyseRegistry(registry)
	if analysis.IsFullyHosted() {
		t.Fatalf("Analysis should have report that the registry is not fully hosted")
	}
	for _, missing := range analysis.MissingHosts {
		if missing == hosted {
			t.Fatalf("Mounted entry should not be reported as missing")
		}
	}
}
