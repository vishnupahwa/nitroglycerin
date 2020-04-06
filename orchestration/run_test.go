package orchestration

import (
	"context"
	"eznft/internal/test/check"
	"github.com/ericchiang/k8s"
	v1 "github.com/ericchiang/k8s/apis/batch/v1"
	"testing"
)

func TestRun(t *testing.T) {
	spec := NFTJob{
		Scenario:        "test",
		Pods:            1,
		Image:           "image",
		MemoryLimit:     "1",
		CPURequest:      "1",
		StartTime:       1,
		OrchestratorUri: "uri",
	}

	jobClient := StubbedJobClient{
		t:            t,
		expectedSpec: spec,
		succeeded:    1,
		failed:       0,
	}

	_, err := Run(&jobClient, spec)
	check.Ok(t, err)
}

func TestRun_fails(t *testing.T) {
	spec := NFTJob{
		Scenario:        "test",
		Pods:            1,
		Image:           "image",
		MemoryLimit:     "1",
		CPURequest:      "1",
		StartTime:       1,
		OrchestratorUri: "uri",
	}

	jobClient := StubbedJobClient{
		t:            t,
		expectedSpec: spec,
		succeeded:    0,
		failed:       1,
	}

	_, err := Run(&jobClient, spec)
	check.Assert(t, err != nil, "Expected an error, got nil")
}

type StubbedJobClient struct {
	t            *testing.T
	deleteCalled bool
	expectedSpec NFTJob
	succeeded    int32
	failed       int32
}

func (s *StubbedJobClient) CreateAndWatch(_ context.Context, spec NFTJob, watcher func(*v1.Job, error, *bool)) {
	check.Equals(s.t, s.expectedSpec, spec)
	job := v1.Job{
		Status: &v1.JobStatus{
			Succeeded: k8s.Int32(s.succeeded),
			Failed:    k8s.Int32(s.failed),
		},
	}
	watcher(&job, nil, k8s.Bool(true))
}

func (s *StubbedJobClient) Delete(NFTJob) {
	s.deleteCalled = true
}
