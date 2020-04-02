package scenario

import (
	"context"
	"eznft/orchestration"
	"eznft/stages"
	"github.com/schollz/progressbar/v3"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Scenario struct {
	Stages         []stages.Stage
	Targets        []vegeta.Target
	TargetModifier func(target vegeta.Target) vegeta.Target
}

func (s *Scenario) Run(ctx context.Context, name string, stream *orchestration.Stream) vegeta.Results {
	run := vegeta.Results{}
	metrics := vegeta.Metrics{}

	attacker := vegeta.NewAttacker()
	targeter := StaticInterceptedTargeter(s.TargetModifier, s.Targets...)
	for i, stage := range s.Stages {
		log.Println("Running stage " + strconv.Itoa(i))
		must(stage.Before())
		for res := range attacker.Attack(targeter, stage.StgPacer, stage.StgDuration, name) {
			if stream != nil {
				must(stream.SendResults(ctx, *res))
			}
			run.Add(res)
			metrics.Add(res)
		}
		must(stage.After())
	}

	run.Close()
	metrics.Close()
	if stream != nil {
		must(stream.Finish(ctx))
	}
	reporter := vegeta.NewTextReporter(&metrics)
	_ = reporter.Report(os.Stdout)
	return run
}

func (s *Scenario) StartProgressBar(ticker *time.Ticker) {
	bar := progressbar.NewOptions(s.TotalTimeSeconds(), progressbar.OptionSetRenderBlankState(true))
	go func() {
		<-ticker.C
		_ = bar.Add(1)
	}()
}

func (s *Scenario) TotalTimeSeconds() int {
	total := 0
	for _, stage := range s.Stages {
		total += int(stage.StgDuration.Seconds())
	}
	return total
}

func must(err error) {
	if err != nil {
		log.Fatalln("Failed to run stage: " + err.Error())
	}
}

// Writes the results in CSV format to the writer
func WriteResults(w io.Writer, results vegeta.Results) {
	csvEncoder := vegeta.NewCSVEncoder(w)
	for _, r := range results {
		err := csvEncoder.Encode(&r)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
