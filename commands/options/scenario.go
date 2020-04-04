package options

import (
	"github.com/spf13/cobra"
	"time"
)

// Run struct contains options regarding the scenario
type Scenario struct {
	StartAt    int64
	StartNext  int
	UploadURI  string
	Multiplier float64
}

func AddStartAtArg(cmd *cobra.Command, r *Scenario) {
	cmd.Flags().Int64VarP(&r.StartAt, "start-at", "s", time.Now().Unix(), "At what Unix seconds to start the scenario. (The number of seconds elapsed since January 1, 1970 UTC)")
}

func AddStartNextArg(cmd *cobra.Command, r *Scenario) {
	cmd.Flags().IntVarP(&r.StartNext, "start-next", "n", 0, "At what next minute increment to start scenario. Determines 'start-at' flag by calculating the Unix time for that exact next minute.")
}

func AddUploadURIArg(cmd *cobra.Command, r *Scenario) {
	cmd.Flags().StringVarP(&r.UploadURI, "upload-uri", "u", "", "URI for streaming results via GRPC")
}

func AddMultiplierArg(cmd *cobra.Command, r *Scenario) {
	cmd.Flags().Float64VarP(&r.Multiplier, "multiplier", "m", 1, "Multiplier for targets for all stages of a scenario")
}
