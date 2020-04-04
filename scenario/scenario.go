package scenario

import (
	"context"
	"eznft/stages"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"log"
	"os"
	"strconv"
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

// Stream represents a stream where results can be consumed.
type Stream interface {
	SendResults(context.Context, vegeta.Result) error
	Close(context.Context) error
}

// Run runs the scenario with the given multiplier, sending all results to the provided stream.
func (s *Scenario) Run(ctx context.Context, name string, multiplier float64, stream Stream) vegeta.Metrics {
	s.stages = s.StagesToBe.Build(multiplier)

	metrics := vegeta.Metrics{}

	attacker := vegeta.NewAttacker()
	targeter := StaticInterceptedTargeter(s.TargetModifier, s.Targets...)
	for i, stage := range s.stages {
		log.Println("Running stage " + strconv.Itoa(i))
		for res := range attacker.Attack(targeter, stage.StgPacer, stage.StgDuration, name) {
			must(stream.SendResults(ctx, *res))
			metrics.Add(res)
		}
	}

	metrics.Close()
	must(stream.Close(ctx))
	reporter := vegeta.NewTextReporter(&metrics)
	_ = reporter.Report(os.Stdout)
	return metrics
}

func must(err error) {
	if err != nil {
		log.Fatalln("Failed to run stage due to stream: " + err.Error())
	}
}
