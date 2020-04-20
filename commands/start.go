package commands

import (
	"context"
	"eznft/commands/options"
	"eznft/definitions"
	"eznft/orchestration"
	"eznft/scenario"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// addStart adds the primary start command to a top level command.
// This is the entrypoint command for starting a controlled application.
func addStart(topLevel *cobra.Command) {
	scenarioOpts := &options.Scenario{}

	startCmd := &cobra.Command{
		Use:   "start <name>",
		Short: "Start an NFT scenario",
		Long: `Start an NFT scenario
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := start(args[0], scenarioOpts)
			if err != nil {
				log.Fatal(err)
			}
		},
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"init", "run"},
	}

	options.AddStartAtArg(startCmd, scenarioOpts)
	options.AddStartNextArg(startCmd, scenarioOpts)
	options.AddUploadURIArg(startCmd, scenarioOpts)
	options.AddMultiplierArg(startCmd, scenarioOpts)
	options.AddTargetOverrideArg(startCmd, scenarioOpts)
	topLevel.AddCommand(startCmd)
}

func start(name string, scenarioOpts *options.Scenario) error {
	now := time.Now()
	ctx := context.Background()
	s, ok := definitions.NFT[name]
	if !ok {
		return fmt.Errorf("scenario %s not found. Possible options: %v", name, getKeys(definitions.NFT))
	}
	if scenarioOpts.StartNext != 0 {
		scenarioOpts.StartAt = orchestration.CalculateStartAt(now, scenarioOpts.StartNext)
	}

	orchestration.WaitForStart(scenarioOpts.StartAt)

	log.Println("Running NFT " + name)

	s.Run(ctx, name, scenarioOpts.Multiplier, scenarioOpts.TargetOverride, stream(ctx, name, scenarioOpts))
	return nil
}

func stream(ctx context.Context, name string, scenarioOpts *options.Scenario) scenario.Stream {
	if scenarioOpts.UploadURI != "" {
		stream, err := orchestration.EstablishStream(ctx, scenarioOpts.UploadURI)
		if err != nil {
			log.Fatalln(err)
		}
		return stream
	}
	return scenario.EstablishCSV(name, time.Now())
}

func getKeys(m map[string]scenario.Scenario) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
