package definitions

import (
	"eznft/scenario"
	"eznft/stages"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"time"
)

func example() scenario.Scenario {
	return scenario.Scenario{
		StagesToBe: stages.Builder().
			// Ramp up to 10 TPS over 5 seconds and sustain for 20 seconds
			RampUpAndSustain(10, 10*time.Second, 10*time.Second),
		Targets: []vegeta.Target{{
			Method: "POST",
			URL:    "http://hits",
			Body:   nil,
			Header: nil,
		}},
	}
}
