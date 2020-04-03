package scenario

import (
	"context"
	"eznft/orchestration"
	"eznft/stages"
	"github.com/schollz/progressbar/v3"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// Scenario defines a single non-functional test scenario to be ran or orchestrated.
// It is made up of 'stages' and targets to be hit.
type Scenario struct {
	// StagesToBe is a builder that returns a list of stages after being modified with a multiplier
	StagesToBe *stages.B
	// Targets specify the endpoints to be hit
	Targets []vegeta.Target
	// TargetModifier modifies the endpoints before being used. Must be safe for concurrent use.
	TargetModifier func(target vegeta.Target) vegeta.Target
	// stages should be populated via building StagesToBe
	stages []stages.Stage
}

// Run runs the scenario with the given multiplier.
// If an orchestration stream is provided, all results will be sent there instead of being written to a file
func (s *Scenario) Run(ctx context.Context, name string, multiplier float64, stream *orchestration.Stream) vegeta.Metrics {
	s.stages = s.StagesToBe.Build(multiplier)

	metrics := vegeta.Metrics{}

	attacker := vegeta.NewAttacker()
	targeter := StaticInterceptedTargeter(s.TargetModifier, s.Targets...)
	for i, stage := range s.stages {
		log.Println("Running stage " + strconv.Itoa(i))
		for res := range attacker.Attack(targeter, stage.StgPacer, stage.StgDuration, name) {
			if stream != nil {
				must(stream.SendResults(ctx, *res))
			}
			metrics.Add(res)
		}
	}

	metrics.Close()
	if stream != nil {
		must(stream.Finish(ctx))
	}
	reporter := vegeta.NewTextReporter(&metrics)
	_ = reporter.Report(os.Stdout)
	return metrics
}

func (s *Scenario) StartProgressBar(ticker *time.Ticker) {
	bar := progressbar.NewOptions(s.TotalTimeSeconds(), progressbar.OptionSetRenderBlankState(true))
	go func() {
		<-ticker.C
		_ = bar.Add(1)
	}()
}

func (s *Scenario) TotalTimeSeconds() int {
	total := 0
	for _, stage := range s.stages {
		total += int(stage.StgDuration.Seconds())
	}
	return total
}

func must(err error) {
	if err != nil {
		log.Fatalln("Failed to run stage due to stream: " + err.Error())
	}
}

// Writes the results in CSV format to the writer
func WriteResults(w io.Writer, results vegeta.Results) {
	csvEncoder := vegeta.NewCSVEncoder(w)
	for _, r := range results {
		err := csvEncoder.Encode(&r)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
