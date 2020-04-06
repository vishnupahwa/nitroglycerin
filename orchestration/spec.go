package orchestration

type NFTJob struct {
	Scenario        string
	Pods            int32
	Image           string
	MemoryLimit     string
	CPURequest      string
	StartTime       int64
	OrchestratorUri string
	Args            []string
}
