package definitions

import (
	"encoding/json"
	"eznft/scenario"
	"eznft/stages"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"log"
	"net/http"
	"time"
)

func quick() scenario.Scenario {
	return scenario.Scenario{
		StagesToBe: stages.Builder().
			// Ramp up to 10 TPS over 5 seconds and sustain for 20 seconds
			RampUpAndSustain(10, 5*time.Second, 20*time.Second).
			// Ramp down to 0 TPS
			RampDown(5 * time.Second),
		Targets: []vegeta.Target{{
			Method: "POST",
			URL:    "https://echo-r2oihcniea-ew.a.run.app/playout/live",
			Body:   createBody(),
			Header: createHeaders(),
		}},
	}
}

func createBody() []byte {
	type body struct {
		ContentId string `json:"contentId"`
	}
	b := body{ContentId: "contentId"}
	marshalled, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	return marshalled
}

func createHeaders() http.Header {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	return headers
}
