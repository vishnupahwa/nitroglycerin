package definitions

import (
	"eznft/scenario"
	"eznft/stages"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"time"
)

func slow() scenario.Scenario {
	return scenario.Scenario{
		Stages: stages.Builder().
			// Ramp up to 100 TPS
			Add(100, 1*time.Minute).
			// Constant 100 TPS
			Add(100, 10*time.Minute).
			// Ramp down to 0 TPS
			Add(0, 10*time.Second).
			Build(),
		Targets: []vegeta.Target{{
			Method: "POST",
			URL:    "https://sps-lite-2-r2oihcniea-ew.a.run.app/playout/live",
			Body:   createBody(),
			Header: createHeaders(),
		}},
	}
}
