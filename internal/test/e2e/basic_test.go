//+build e2e

// e2e package has end to end tests for eznft.
// As a prerequisite, kind must be running, with the Kubernetes context named kind-eznft.
// Run the ./hack/e2e-setup.sh for this setup.
package e2e

import (
	"context"
	"eznft/internal/test/check"
	"fmt"
	"github.com/ericchiang/k8s"
	batchv1 "github.com/ericchiang/k8s/apis/batch/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

// loadClient parses a kubeconfig from a file and returns a Kubernetes
// client. It does not support extensions or client auth providers.
func loadClient(kubeconfigPath string) (*k8s.Client, error) {
	data, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}

	// Unmarshal YAML into a Kubernetes config object.
	var config k8s.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return k8s.NewClient(&config)
}

func eznftJob() *batchv1.Job {
	return &batchv1.Job{
		Metadata: &metav1.ObjectMeta{
			Name:      k8s.String("eznft"),
			Namespace: k8s.String(clusterNamespace),
		}}
}

const clusterContext = "kind-eznft"
const clusterNamespace = "default"

func TestEndToEnd(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	kubeConfig := homeDir + "/.kube/config"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelFunc()

	command := exec.CommandContext(ctx, "bash", "-c", "kustomize build ./testdata/resources/kind-test/ | kubectl --context "+clusterContext+" apply -f -")
	err := command.Run()
	check.Ok(t, err)

	client, err := loadClient(kubeConfig)
	check.Ok(t, err)
	watch, err := client.Watch(ctx, clusterNamespace, eznftJob())
	check.Ok(t, err)
	watching := true
	for watching {
		j := new(batchv1.Job)
		_, err := watch.Next(j)
		check.Ok(t, err)
		if *j.Status.Succeeded > 0 {
			watching = false
		}
		if *j.Status.Failed > 0 {
			watching = false
			t.Errorf("Expected job to succeed.")
		}
	}

	command = exec.CommandContext(ctx, "bash", "-c", "kubectl --context "+clusterContext+" run -i --rm --quiet --restart=Never alpine --image=alpine -- wget -qO - http://hits/count")
	output, err := command.CombinedOutput()
	check.Ok(t, err)
	count, err := strconv.Atoi(string(output))
	check.Ok(t, err)
	check.Equals(t, 160, count)
}
