package orchestration

import (
	"context"
	"encoding/json"
	"eznft/internal/test/check"
	"eznft/orchestration/proto"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func TestEstablishStream(t *testing.T) {

}

func TestStream_SendResults(t *testing.T) {
	ctx := context.Background()
	client := StubbedOrchestratorClient{}
	s := Stream{clc: &StubbedCloseableConn{}, sender: &client}
	defer s.Close(ctx)
	result := vegeta.Result{
		Attack:    "test",
		Seq:       0,
		Code:      200,
		Timestamp: time.Time{},
		Body:      nil,
		Method:    "GET",
		URL:       "localhost:8080",
	}

	err := s.SendResults(ctx, result)
	check.Ok(t, err)

	marshaled, _ := json.Marshal(result)
	check.Equals(t, marshaled, client.result)
}

func TestStream_Close(t *testing.T) {
	ctx := context.Background()
	client := StubbedOrchestratorClient{}
	conn := StubbedCloseableConn{}
	s := Stream{clc: &conn, sender: &client}
	err := s.Close(ctx)
	check.Ok(t, err)
	check.Equals(t, true, client.closed)
	check.Equals(t, true, conn.closed)
}

func TestUploadingToServer(t *testing.T) {
	orch, cancelFunc := startUploadServer()
	defer cancelFunc()
	ctx := context.Background()
	stream, err := EstablishStream(ctx, "localhost:8080")
	check.Ok(t, err)
	result := vegeta.Result{
		Attack:    "test",
		Seq:       0,
		Code:      200,
		Timestamp: time.Time{},
		Body:      nil,
		Method:    "GET",
		URL:       "localhost:8080",
	}
	err = stream.SendResults(ctx, result)
	check.Ok(t, err)
	check.Ok(t, stream.Close(ctx))

	orch.Close()
	check.Equals(t, uint64(1), orch.metrics.Requests)
}

type StubbedOrchestratorClient struct {
	result []byte
	closed bool
}

func (s *StubbedOrchestratorClient) Send(r *orchestration.Result) error {
	s.result = r.Json
	return nil
}

func (s *StubbedOrchestratorClient) CloseAndRecv() (*orchestration.Complete, error) {
	s.closed = true
	return nil, nil
}

func (s *StubbedOrchestratorClient) Header() (metadata.MD, error) {
	panic("implement me")
}

func (s StubbedOrchestratorClient) Trailer() metadata.MD {
	panic("implement me")
}

func (s StubbedOrchestratorClient) CloseSend() error {
	panic("implement me")
}

func (s StubbedOrchestratorClient) Context() context.Context {
	panic("implement me")
}

func (s StubbedOrchestratorClient) SendMsg(m interface{}) error {
	panic("implement me")
}

func (s StubbedOrchestratorClient) RecvMsg(m interface{}) error {
	panic("implement me")
}

type StubbedCloseableConn struct {
	closed bool
}

func (s *StubbedCloseableConn) Close() error {
	s.closed = true
	return nil
}
