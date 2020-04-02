// Package stages provides the stage type which makeup scenarios.
package stages

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"time"
)

/* Stage is a section of a scenario such as the ramp up, constant load or the ramp down.
Stage represents the TPS approaching a target for a time period.
For example the initial ramp up would be a separate stage to the constant load that follows for a peak load test. */
type Stage struct {
	Before func() error
	// Target specifies the transaction per second target of this stage
	Target int
	// StgDuration is the total duration of this stage.
	StgDuration time.Duration
	// Pacer determines the slope for this stage
	StgPacer vegeta.Pacer
	After    func() error
}

// NewStage creates a stage with the provided target transactions per second,
// aiming to either ramp to that rate or sustain that rate for a given stgDuration,
// depending on the previous stage provided.
func NewStage(target int, duration time.Duration, prev Stage) Stage {
	pacer := pacerFrom(Stage{Target: target, StgDuration: duration}, prev)
	return Stage{Before: NoOperation(), Target: target, StgDuration: duration, StgPacer: pacer, After: NoOperation()}
}

func pacerFrom(current, previous Stage) vegeta.Pacer {
	if current.Target != previous.Target {
		return current.rampFrom(previous)
	}
	return vegeta.ConstantPacer{
		Freq: current.Target,
		Per:  1 * time.Second,
	}
}

func (s *Stage) rampFrom(previous Stage) vegeta.Pacer {
	return vegeta.LinearPacer{
		StartAt: vegeta.Rate{Freq: previous.Target, Per: 1 * time.Second},
		Slope:   calculateSlope(previous.Target, s.Target, s.StgDuration),
	}
}

func calculateSlope(base int, target int, duration time.Duration) float64 {
	return float64(target-base) / (duration.Seconds())
}

func NoOperation() func() error {
	return func() error {
		// Do nothing
		return nil
	}
}
