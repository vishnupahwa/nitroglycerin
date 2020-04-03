package commands

import (
	"errors"
	"eznft/commands/options"
	"eznft/definitions"
	"eznft/orchestration"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// addOrchestrate adds the orchestrate command
func addOrchestrate(topLevel *cobra.Command) {
	orchestrationOpts := &options.Orchestration{}

	orchCmd := &cobra.Command{
		Use:   "orchestrate <name>",
		Short: "Orchestrate an NFT scenario",
		Long: `Orchestrate an NFT scenario, creating distributed load through 
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := orchestrate(args[0], orchestrationOpts)
			if err != nil {
				log.Fatal(err)
			}
		},
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"init", "run"},
	}

	options.AddPodsArg(orchCmd, orchestrationOpts)
	options.AddImageArg(orchCmd, orchestrationOpts)
	options.AddSelfURIArg(orchCmd, orchestrationOpts)
	options.AddCPURequestArg(orchCmd, orchestrationOpts)
	options.AddMemoryLimitsArg(orchCmd, orchestrationOpts)
	topLevel.AddCommand(orchCmd)
}

func orchestrate(name string, opts *options.Orchestration) error {
	_, ok := definitions.NFT[name]
	if !ok {
		return errors.New(fmt.Sprintf("Scenario %s not found. Possible options: %v", name, getKeys(definitions.NFT)))
	}
	log.Println("Orchestrating scenario " + name)
	now := time.Now()
	startAt := orchestration.CalculateStartAt(now, 2)
	log.Printf("Scenario execution %s starting at %v", name, time.Unix(startAt, 0))
	spec := orchestration.NFTJob{
		Scenario:        name,
		Pods:            opts.Pods,
		Image:           opts.Image,
		CPURequest:      opts.CPURequests,
		MemoryLimit:     opts.MemoryLimits,
		StartTime:       startAt,
		OrchestratorUri: opts.SelfURI,
	}
	orchestration.Run(spec)
	return nil
}
