package request

import (
	"github.com/giantswarm/kubernetesdclient/service/creator/request/aws"
)

// Worker configures the Kubernetes worker nodes.
type Worker struct {
	AWS     aws.Worker        `json:"aws"`
	CPU     CPU               `json:"cpu"`
	ID      string            `json:"id"`
	Labels  map[string]string `json:"labels"`
	Memory  Memory            `json:"memory"`
	Storage Storage           `json:"storage"`
}

// DefaultWorker provides a default worker configuration by best effort.
func DefaultWorker() Worker {
	return Worker{
		AWS:     aws.DefaultWorker(),
		CPU:     DefaultCPU(),
		ID:      "",
		Labels:  map[string]string{},
		Memory:  DefaultMemory(),
		Storage: DefaultStorage(),
	}
}
