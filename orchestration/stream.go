package orchestration

import (
	"context"
	"encoding/json"
	pb "eznft/orchestration/proto"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"google.golang.org/grpc"
	"os"
)

type Stream struct {
	clc    CloseableConn
	sender pb.Orchestrator_SendResultsClient
}

type CloseableConn interface {
	Close() error
}

func EstablishStream(ctx context.Context, URI string) (*Stream, error) {
	conn, err := grpc.Dial(URI, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	sender, err := pb.NewOrchestratorClient(conn).SendResults(ctx)
	if err != nil {
		return nil, err
	}
	return &Stream{clc: conn, sender: sender}, nil
}

func (c *Stream) SendResults(ctx context.Context, result vegeta.Result) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return c.sender.Send(&pb.Result{
		Json:     b,
		Hostname: hostname,
	})
}

func (c *Stream) Close(ctx context.Context) error {
	defer c.clc.Close()
	_, err := c.sender.CloseAndRecv()
	if err != nil {
		return err
	}
	return nil
}
