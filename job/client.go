package job

import (
	"context"
	"eznft/orchestration"
	"fmt"
	"github.com/ericchiang/k8s"
	batchv1 "github.com/ericchiang/k8s/apis/batch/v1"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/ericchiang/k8s/apis/resource"
	"log"
	"os"
	"strconv"
)

type Client struct {
	Client *k8s.Client
}

func CreateClient() *Client {
	c, err := k8s.NewInClusterClient()
	if err != nil {
		log.Fatal(err)
	}
	return &Client{Client: c}
}

func (c *Client) Delete(spec orchestration.NFTJob) {
	// Delete existing job if it exists
	jobRef := &batchv1.Job{Metadata: &metav1.ObjectMeta{Name: k8s.String(spec.Scenario), Namespace: k8s.String(c.Client.Namespace)}}
	_ = c.Client.Delete(context.Background(), jobRef)
}

func (c *Client) jobFrom(spec orchestration.NFTJob) *batchv1.Job {
	pullPolicy := k8s.String("Always")
	if _, isMinikube := os.LookupEnv("MINIKUBE"); isMinikube {
		pullPolicy = k8s.String("Never")
	}
	return &batchv1.Job{
		Metadata: &metav1.ObjectMeta{
			Name:      k8s.String(spec.Scenario),
			Namespace: k8s.String(c.Client.Namespace),
		},
		Spec: &batchv1.JobSpec{
			Parallelism:  k8s.Int32(spec.Pods),
			Completions:  k8s.Int32(spec.Pods),
			BackoffLimit: k8s.Int32(0),
			Template: &corev1.PodTemplateSpec{
				Spec: &corev1.PodSpec{
					Containers: []*corev1.Container{
						{
							Name:  k8s.String(spec.Scenario),
							Image: k8s.String(spec.Image),
							Args: []string{
								"start",
								spec.Scenario,
								"--start-at",
								strconv.FormatInt(spec.StartTime, 10),
								"--upload-uri",
								spec.OrchestratorUri,
								"--multiplier",
								fmt.Sprintf("%.2f", 1/float64(spec.Pods)),
							},
							ImagePullPolicy: pullPolicy,
							Resources: &corev1.ResourceRequirements{
								Limits: map[string]*resource.Quantity{
									"memory": {String_: k8s.String(spec.MemoryLimit)},
								},
								Requests: map[string]*resource.Quantity{
									"cpu": {String_: k8s.String(spec.CPURequest)},
								},
							},
							Env: []*corev1.EnvVar{
								{
									Name:  k8s.String("ORCHESTRATION_URI"),
									Value: k8s.String(spec.OrchestratorUri),
								},
							},
						},
					},
					RestartPolicy: k8s.String("Never"),
				},
			},
		},
	}
}

func (c *Client) create(ctx context.Context, spec orchestration.NFTJob) {
	err := c.Client.Create(ctx, c.jobFrom(spec))
	if err != nil {
		log.Fatalf("Could not create Job: %v\n", err)
	}
}

func (c *Client) CreateAndWatch(ctx context.Context, spec orchestration.NFTJob, watcher func(j *batchv1.Job, err error, running *bool)) {
	c.create(ctx, spec)
	watch, err := c.Client.Watch(ctx, c.Client.Namespace, c.jobFrom(spec))
	if err != nil {
		log.Fatalf("Could not watch Job: %v\n", err)
	}
	running := true
	for running {
		j := new(batchv1.Job)
		_, err := watch.Next(j)
		watcher(j, err, &running)
	}
}
