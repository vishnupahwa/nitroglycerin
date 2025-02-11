package definitions

import (
	"eznft/scenario"
	"eznft/stages"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"time"
)

func slower() scenario.Scenario {
	return scenario.Scenario{
		StagesToBe: stages.Builder().
			// Ramp up to 100 TPS
			RampUpAndSustain(15000, 10*time.Minute, 1*time.Second).
			// Ramp down to 0 TPS
			RampDown(1 * time.Second),
		Targets: []vegeta.Target{{
			Method: "POST",
			URL:    "https://echo-r2oihcniea-ew.a.run.app/playout/live",
			Body:   createBody(),
			Header: createHeaders(),
		}},
	}
}
