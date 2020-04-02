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
		Stages: stages.Builder().
			// Ramp up to 100 TPS
			Add(10, 5*time.Second).
			// Constant 100 TPS
			Add(10, 20*time.Second).
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
