package gorch

type RegistryAnalysis struct {
	MissingHosts []*OperationEntry
	Registry     *OperationsRegistry
}

func AnalyseRegistry(registry *OperationsRegistry) *RegistryAnalysis {
	analysis := &RegistryAnalysis{
		Registry: registry,
	}
	analysis.analyze()
	return analysis
}

func (analysis *RegistryAnalysis) analyze() {
	for _, op := range analysis.Registry.Operations() {
		if !op.IsHosted() {
			analysis.MissingHosts = append(analysis.MissingHosts, op)
		}
	}
}

func (analysis *RegistryAnalysis) IsFullyHosted() bool {
	return len(analysis.MissingHosts) == 0
}
