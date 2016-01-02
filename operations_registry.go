package gorch

import (
	"log"
	"sync"
)

type OperationsRegistry struct {
	operations map[string]*OperationEntry
	init       sync.Once
}

func (registry *OperationsRegistry) initialize() {
	registry.init.Do(func() {
		registry.operations = make(map[string]*OperationEntry)
	})
}

func (registry *OperationsRegistry) op(name string) *OperationEntry {
	registry.initialize()
	op, ok := registry.operations[name]
	if !ok {
		return nil
	}
	return op
}

func (registry *OperationsRegistry) Exists(name string) bool {
	op := registry.op(name)
	return registry.exists(op)
}

func (registry *OperationsRegistry) exists(op *OperationEntry) bool {
	return op != nil
}

func (registry *OperationsRegistry) Register(name string) *OperationEntry {
	if registry.Exists(name) {
		log.Fatalf("Operation %s is already registered", name)
	}
	entry := &OperationEntry{
		name: name,
	}
	registry.addEntry(entry)
	return entry
}

func (registry *OperationsRegistry) addEntry(entry *OperationEntry) {
	registry.operations[entry.Name()] = entry
}

func (registry *OperationsRegistry) Entry(name string) *OperationEntry {
	op := registry.op(name)
	return op
}

func (registry *OperationsRegistry) Operations() []*OperationEntry {
	registry.initialize()
	ops := make([]*OperationEntry, 0, len(registry.operations))
	for _, op := range registry.operations {
		ops = append(ops, op)
	}
	return ops
}

func (registry *OperationsRegistry) Merge(source *OperationsRegistry) {
	registry.initialize()
	for _, op := range source.Operations() {
		name := op.Name()
		if registry.Exists(name) {
			log.Fatalf("Registry merge failed, operation %s is already registered", name)
		}
		entry := &OperationEntry{}
		entry.merge(op)
		registry.addEntry(entry)
	}
}

func (registry *OperationsRegistry) Analyze() *RegistryAnalysis {
	return AnalyseRegistry(registry)
}
