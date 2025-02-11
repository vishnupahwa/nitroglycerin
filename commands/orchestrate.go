package commands

import (
	"eznft/commands/options"
	"eznft/definitions"
	"eznft/job"
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
		Long: `Orchestrate an NFT scenario, creating distributed load through a Kubernetes Job.

Creates a Kubernetes job which spins up pods for load testing the specified scenario. 
Starts up a GRPC server which the pods stream their load results to. After all pods have
completed, the results are sorted and used to calculate an overall report.

This command itself must be run as a Kubernetes job. 
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := orchestrate(args[0], orchestrationOpts)
			if err != nil {
				log.Fatal(err)
			}
		},
		Args: cobra.ExactArgs(1),
	}

	options.AddPodsArg(orchCmd, orchestrationOpts)
	options.AddImageArg(orchCmd, orchestrationOpts)
	options.AddSelfURIArg(orchCmd, orchestrationOpts)
	options.AddCPURequestArg(orchCmd, orchestrationOpts)
	options.AddMemoryLimitsArg(orchCmd, orchestrationOpts)
	options.AddForwardedArgsArg(orchCmd, orchestrationOpts)
	topLevel.AddCommand(orchCmd)
}

func orchestrate(name string, opts *options.Orchestration) error {
	s, ok := definitions.NFT[name]
	if !ok {
		return fmt.Errorf("scenario %s not found. Possible options: %v", name, getKeys(definitions.NFT))
	}
	log.Println("Orchestrating scenario " + name)
	var _ = s.StagesToBe.Build(1 / float64(opts.Pods)) // Check scenario successfully builds

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
		OrchestratorURI: opts.SelfURI,
		Args:            opts.Args,
	}
	_, err := orchestration.Run(job.CreateClient(), spec)
	return err
}
