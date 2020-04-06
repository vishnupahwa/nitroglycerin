package orchestration

import (
	"context"
	"errors"
	batchv1 "github.com/ericchiang/k8s/apis/batch/v1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"log"
	"os"
)

type client interface {
	CreateAndWatch(context.Context, NFTJob, func(*batchv1.Job, error, *bool))
	Delete(NFTJob)
}

func Run(client client, spec NFTJob) (vegeta.Metrics, error) {
	client.Delete(spec)

	ctx := context.Background()
	orchestrator, cancelFunc := startUploadServer()
	defer cancelFunc()
	succeeded := false

	client.CreateAndWatch(ctx, spec, func(j *batchv1.Job, err error, running *bool) {
		if err != nil {
			log.Println("Error watching jobs: ", err)
			*running = false
		}
		if j.Status != nil && *j.Status.Failed > 0 {
			log.Println("Pod failed")
			*running = false
		}
		if j.Status != nil && *j.Status.Succeeded == spec.Pods {
			log.Println("All orchestrated pods succeeded!")
			*running = false
			succeeded = true
		}
	})

	if !succeeded {
		log.Println("Exiting due to failure")
		return vegeta.Metrics{}, errors.New("exiting due to failure")
	}

	orchestrator.Close()
	err := vegeta.NewTextReporter(orchestrator.metrics).Report(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
	return *orchestrator.metrics, nil

}
