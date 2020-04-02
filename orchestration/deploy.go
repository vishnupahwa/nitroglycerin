package orchestration

import (
	"context"
	pb "eznft/orchestration/proto"
	"fmt"
	"github.com/ericchiang/k8s"
	batchv1 "github.com/ericchiang/k8s/apis/batch/v1"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/ericchiang/k8s/apis/resource"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
)

// TODO break up - too big
func Run(spec NFTJob) (vegeta.Metrics, vegeta.Results) {
	client, err := k8s.NewInClusterClient()
	if err != nil {
		log.Fatal(err)
	}

	// Delete existing job if it exists
	jobRef := &batchv1.Job{Metadata: &metav1.ObjectMeta{Name: k8s.String(spec.Scenario), Namespace: k8s.String(client.Namespace)}}
	_ = client.Delete(context.Background(), jobRef)

	pullPolicy := k8s.String("IfNotPresent")
	if _, isMinikube := os.LookupEnv("MINIKUBE"); isMinikube {
		pullPolicy = k8s.String("Never")
	}
	job := &batchv1.Job{
		Metadata: &metav1.ObjectMeta{
			Name:      k8s.String(spec.Scenario),
			Namespace: k8s.String(client.Namespace),
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
	err = client.Create(context.Background(), job)
	if err != nil {
		log.Fatalf("Could not create Job: %v\n", err)
	}
	ctx := context.Background()
	watcher, err := client.Watch(ctx, client.Namespace, jobRef)
	if err != nil {
		log.Fatalf("Could not watch Job: %v\n", err)
	}
	defer watcher.Close()

	orchestrator, cancelFunc := startUploadServer()
	defer cancelFunc()
	running := true
	succeeded := false
	for running {
		j := new(batchv1.Job)
		_, err := watcher.Next(j)
		if err != nil {
			log.Println("Error watching jobs: ", err)
			running = false
		}
		if *j.Status.Failed > 0 {
			log.Println("Pod failed")
			running = false
		}
		if *j.Status.Succeeded == spec.Pods {
			log.Println("All orchestrated pods succeeded!")
			running = false
			succeeded = true
		}
	}
	if !succeeded {
		log.Println("Exiting due to failure")
		os.Exit(1)
	}

	orchestrator.Close()
	err = vegeta.NewTextReporter(orchestrator.metrics).Report(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
	return *orchestrator.metrics, nil
}

func startUploadServer() (*Orchestrator, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	grpcServer := grpc.NewServer()
	orchestrator := NewOrchestrator()
	pb.RegisterOrchestratorServer(grpcServer, orchestrator)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalln("Could not start listening on 8080: " + err.Error())
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("grpcServer: Serve() error: %s", err)
		}
	}()
	go func() {
		<-ctx.Done()
		println("Stopping GRPC Server")
		grpcServer.GracefulStop()
	}()
	return orchestrator, cancelFunc
}
