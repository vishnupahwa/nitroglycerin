package options

import (
	"github.com/spf13/cobra"
)

// Run struct contains options regarding the scenario
type Orchestration struct {
	Pods         int32
	Image        string
	CPURequests  string
	MemoryLimits string
	SelfURI      string
	Args         []string
}

func AddPodsArg(cmd *cobra.Command, r *Orchestration) {
	cmd.Flags().Int32VarP(&r.Pods, "pods", "p", 5, "Number of pods to distribute load over")
}

func AddImageArg(cmd *cobra.Command, r *Orchestration) {
	cmd.Flags().StringVarP(&r.Image, "image", "i", "eznft:latest", "Image to use for orchestrated Jobs")
}

func AddSelfURIArg(cmd *cobra.Command, r *Orchestration) {
	cmd.Flags().StringVarP(&r.SelfURI, "uri", "u", "eznft:8080", "URI of what is calling eznft for uploading results")
}

func AddCPURequestArg(cmd *cobra.Command, r *Orchestration) {
	cmd.Flags().StringVarP(&r.CPURequests, "cpu", "c", "500m", "CPU resource request of pods")
}

func AddMemoryLimitsArg(cmd *cobra.Command, r *Orchestration) {
	cmd.Flags().StringVarP(&r.MemoryLimits, "memory", "m", "0.5Gi", "Memory resource limits of pods")
}

func AddForwardedArgsArg(cmd *cobra.Command, r *Orchestration) {
	cmd.Flags().StringSliceVarP(&r.Args, "args", "a", nil, "Args to forward to orchestrated pods")
}
