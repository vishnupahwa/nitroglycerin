package definitions

import (
	"eznft/scenario"
	"eznft/stages"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"time"
)

// TODO BUG chaining ramp up stages doesn't work as it uses to the steady up pacer. The old linear pacer is better suited.
func spike() scenario.Scenario {
	return scenario.Scenario{
		StagesToBe: stages.Builder().
			RampUpAndSustain(5000, 3*time.Minute, 1*time.Second).
			RampUpAndSustain(1000, 1*time.Minute, 1*time.Second).
			RampUpAndSustain(4000, 1*time.Minute, 1*time.Second).
			RampUpAndSustain(100, 1*time.Minute, 1*time.Second).
			RampUpAndSustain(3000, 1*time.Minute, 1*time.Second).
			RampDown(1 * time.Second),
		Targets: []vegeta.Target{{
			Method: "POST",
			URL:    "https://example.com",
			Body:   createBody(),
			Header: createHeaders(),
		}},
	}
}
