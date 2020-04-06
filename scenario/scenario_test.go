package scenario

import (
	"context"
	"eznft/internal/test/check"
	"eznft/stages"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testStream struct {
	results []vegeta.Result
	closed  bool
}

func (t *testStream) SendResults(_ context.Context, r vegeta.Result) error {
	t.results = append(t.results, r)
	return nil
}

func (t *testStream) Close(context.Context) error {
	t.closed = true
	return nil
}

func TestScenario_Run(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
	defer server.Close()
	s := Scenario{
		StagesToBe: stages.Builder().AddFixedStage(stages.Stage{
			Target:      1,
			StgDuration: 1 * time.Second,
			StgPacer:    vegeta.ConstantPacer{Freq: 1, Per: time.Second},
		}),
		Targets: []vegeta.Target{{Method: "GET", URL: server.URL}},
	}
	stream := &testStream{}
	run := s.Run(context.Background(), "test", 1, "", stream)
	check.Equals(t, 1.00, run.Success)
	check.Equals(t, 1, len(stream.results))
}
