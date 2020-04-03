package commands

import (
	"context"
	"errors"
	"eznft/commands/options"
	"eznft/definitions"
	"eznft/orchestration"
	"eznft/scenario"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
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
	topLevel.AddCommand(startCmd)
}

func start(name string, scenarioOpts *options.Scenario) error {
	now := time.Now()
	ctx := context.Background()
	s, ok := definitions.NFT[name]
	if !ok {
		return errors.New(fmt.Sprintf("Scenario %s not found. Possible options: %v", name, getKeys(definitions.NFT)))
	}
	if scenarioOpts.StartNext != 0 {
		scenarioOpts.StartAt = orchestration.CalculateStartAt(now, scenarioOpts.StartNext)
	}

	orchestration.WaitForStart(scenarioOpts.StartAt)

	log.Println("Running NFT " + name)

	results := s.Run(ctx, name, scenarioOpts.Multiplier, stream(scenarioOpts, ctx))
	file, err := os.Create("nft-results-" + name + "-" + strconv.FormatInt(now.Unix(), 10) + ".csv")
	if err != nil {
		log.Fatalln(err)
	}
	scenario.WriteResults(file, results)
	return nil
}

func stream(scenarioOpts *options.Scenario, ctx context.Context) *orchestration.Stream {
	if scenarioOpts.UploadURI != "" {
		stream, err := orchestration.EstablishStream(scenarioOpts.UploadURI, ctx)
		if err != nil {
			log.Fatalln(err)
		}
		return stream
	}
	return nil
}

func getKeys(m map[string]scenario.Scenario) []string {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
