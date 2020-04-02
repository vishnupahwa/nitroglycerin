package orchestration

import (
	"encoding/json"
	pb "eznft/orchestration/proto"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"io"
	"log"
	"sync/atomic"
)

//go:generate protoc -I ./proto/ ./proto/orchestration.proto --go_out=plugins=grpc:./proto

type Orchestrator struct {
	results chan vegeta.Result
	// Should not be used concurrently
	metrics *vegeta.Metrics
	count   int32
}

// NewOrchestrator creates a new orchestrator.
func NewOrchestrator() *Orchestrator {
	o := &Orchestrator{results: make(chan vegeta.Result), metrics: &vegeta.Metrics{}}
	go func() {
		for r := range o.results {
			o.metrics.Add(&r)
		}
	}()
	return o
}

func (o *Orchestrator) SendResults(s pb.Orchestrator_SendResultsServer) error {
	for {
		result, err := s.Recv()
		if err == io.EOF {
			atomic.AddInt32(&o.count, 1)
			log.Printf("Received all results from pod! Count: %d", o.count)
			return s.SendAndClose(&pb.Complete{Done: true})
		}
		if err != nil {
			return err
		}
		var vr vegeta.Result
		err = json.Unmarshal(result.Json, &vr)
		if err != nil {
			return err
		}

		o.results <- vr
	}
}

func (o *Orchestrator) Close() {
	o.metrics.Close()
	close(o.results)
}
