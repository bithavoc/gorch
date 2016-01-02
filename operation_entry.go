package gorch

import (
	"log"
)

type OperationEntry struct {
	name    string
	handler OperationHandler
}

func (entry *OperationEntry) Name() string {
	return entry.name
}

func (entry *OperationEntry) Handler() OperationHandler {
	return entry.handler
}

func (entry *OperationEntry) IsHosted() bool {
	return entry.handler != nil
}

func (entry *OperationEntry) Host(handler OperationHandler) {
	if entry.IsHosted() {
		log.Fatalf("Operation %s is already hosted", entry.Name())
	}
	entry.handler = handler
}

func (entry *OperationEntry) merge(source *OperationEntry) {
	entry.name = source.Name()
	entry.Host(source.Handler())
}
