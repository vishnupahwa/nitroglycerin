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

func scalability() scenario.Scenario {
	return scenario.Scenario{
		StagesToBe: stages.Builder().
			// Ramp up to 15000 TPS
			RampUpAndSustain(15000, 5*time.Minute, 10*time.Minute).
			// Ramp down to 0 TPS
			RampDown(1 * time.Minute),
		Targets: []vegeta.Target{{
			Method: "POST",
			URL:    "https://echo-r2oihcniea-ew.a.run.app/playout/live",
			Body:   creatLiteBody(),
			Header: createLiteHeaders(),
		}},
	}
}

func creatLiteBody() []byte {
	type body struct {
		ContentID string `json:"contentId"`
	}
	b := body{ContentID: "contentId"}
	marshalled, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	return marshalled
}

func createLiteHeaders() http.Header {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("X-Sky-Signature", "hmacSignature")
	headers.Add("X-SkyOTT-Provider", "provider")
	headers.Add("X-SkyOTT-Territory", "territory")
	return headers
}
