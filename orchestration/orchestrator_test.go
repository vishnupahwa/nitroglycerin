package orchestration

import (
	"context"
	"encoding/json"
	"eznft/internal/test/check"
	"eznft/orchestration/proto"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"google.golang.org/grpc/metadata"
	"io"
	"testing"
	"time"
)

func TestNewOrchestrator(t *testing.T) {
	orchestrator := NewOrchestrator()
	defer orchestrator.Close()
	result := vegeta.Result{
		Attack:    "test",
		Seq:       0,
		Code:      200,
		Timestamp: time.Time{},
		Body:      nil,
		Method:    "GET",
		URL:       "localhost:8080",
	}
	orchestrator.results <- result
	// TODO replace with timeout assertion
	time.Sleep(1 * time.Second)
	check.Equals(t, uint64(1), orchestrator.metrics.Requests)
}

func TestOrchestrator_SendResults(t *testing.T) {
	orchestrator := &Orchestrator{results: make(chan vegeta.Result), metrics: &vegeta.Metrics{}}
	want := vegeta.Result{
		Attack:    "test",
		Seq:       0,
		Code:      200,
		Timestamp: time.Time{},
		Body:      nil,
		Method:    "GET",
		URL:       "localhost:8080",
	}
	stubbedSendResultsServer := &StubbedSendResultsServer{
		result: want,
	}
	go func() {
		err := orchestrator.SendResults(stubbedSendResultsServer)
		check.Ok(t, err)
	}()

	stubbedSendResultsServer.send = true
	got := <-orchestrator.results
	check.Equals(t, want, got)
	stubbedSendResultsServer.EOF = true
	// TODO replace with timeout assertion
	time.Sleep(1 * time.Second)
	check.Equals(t, int32(1), orchestrator.count)
	orchestrator.Close()
}

type StubbedSendResultsServer struct {
	send   bool
	EOF    bool
	result vegeta.Result
}

func (s *StubbedSendResultsServer) SendAndClose(*orchestration.Complete) error {
	return nil
}

func (s *StubbedSendResultsServer) Recv() (*orchestration.Result, error) {
	for {
		if s.send {
			b, _ := json.Marshal(s.result)
			s.send = false
			return &orchestration.Result{Json: b, Hostname: "hostname"}, nil
		}
		if s.EOF {
			return nil, io.EOF
		}
	}
}

func (s *StubbedSendResultsServer) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (s *StubbedSendResultsServer) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (s *StubbedSendResultsServer) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (s *StubbedSendResultsServer) Context() context.Context {
	panic("implement me")
}

func (s *StubbedSendResultsServer) SendMsg(m interface{}) error {
	panic("implement me")
}

func (s *StubbedSendResultsServer) RecvMsg(m interface{}) error {
	panic("implement me")
}
