// Package stages provides the stage type which makeup scenarios.
package stages

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"log"
	"time"
)

/* Stage is a section of a scenario such as the ramp up, constant load or the ramp down.
Stage represents the TPS approaching a target for a time period.
For example the initial ramp up would be a separate stage to the constant load that follows for a peak load test. */
type Stage struct {
	// Target specifies the transaction per second target of this stage
	Target int
	// StgDuration is the total duration of this stage.
	StgDuration time.Duration
	// Pacer determines the slope for this stage
	StgPacer vegeta.Pacer
}

// StageToBe returns a stage after it has been modifier with a target multiplier
type StageToBe func(multiplier float64, prev Stage) Stage

// NewRampingStage creates a stage with the provided target transactions per second
// aiming to steadyily ramp to that rate from 0 over the ramp duration and then sustain
// that rate for a given sustain duration
func NewRampingStage(target int, ramp, sustain time.Duration) Stage {
	pacer, err := NewSteadyUp(
		vegeta.Rate{
			Freq: 1,
			Per:  time.Second,
		},
		vegeta.Rate{
			Freq: target,
			Per:  time.Second,
		},
		ramp)
	if err != nil {
		log.Fatal("Failed to create stage: " + err.Error())
	}
	return Stage{
		Target:      target,
		StgDuration: ramp + sustain,
		StgPacer:    pacer,
	}
}

func NewRampDownStage(total time.Duration, prev Stage) Stage {
	return Stage{
		Target:      0,
		StgDuration: total,
		StgPacer: vegeta.LinearPacer{
			StartAt: vegeta.Rate{
				Freq: prev.Target, Per: time.Second,
			},
			Slope: float64(-prev.Target) / (total.Seconds()),
		},
	}
}
