package options

import (
	"github.com/spf13/cobra"
)

// Run struct contains options regarding the scenario
type Orchestration struct {
	Pods    int32
	Image   string
	SelfURI string

	Store bool
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

func AddStoreArg(cmd *cobra.Command, r *Orchestration) {
	cmd.Flags().BoolVarP(&r.Store, "store", "s", false, "Store all results instead of just metrics")
}
