package exporter

import (
	"context"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	microerror "github.com/giantswarm/microkit/error"
)

type pod struct {
	Name              string
	Status            string
	ContainerStatuses []v1.ContainerStatus
}
type namespace struct {
	Name string
	Pods []pod
}
type Response struct {
	Namespaces []namespace
	LastUpdate time.Time
}

type Config struct {
	Excludes []string
}

// New creates a new configured version service.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Excludes == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Excludes must not be empty")
	}
	newService := &Service{
		Config: config,
	}

	return newService, nil
}

// Service implements the version service interface.
type Service struct {
	Config
}

func (s *Service) Get(ctx context.Context) (*Response, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Response{
		Namespaces: getNamespaces(clientset, s.Config),
		LastUpdate: time.Now(),
	}, nil
}

func getNamespaces(c *kubernetes.Clientset, config Config) []namespace {
	// creates the in-cluster config
	ns, err := c.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	s := []namespace{}
	for _, n := range ns.Items {

		if !isExcluded(config.Excludes, n.Name) {
			s = append(s, namespace{
				Name: n.Name,
				Pods: getPods(c, n.Name),
			})
		}
	}

	return s
}

func isExcluded(excludes []string, ns string) bool {
	for _, e := range excludes {
		if e == ns {
			return true
		}
	}

	return false
}

func getPods(c *kubernetes.Clientset, ns string) []pod {
	// creates the in-cluster config
	pods, err := c.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	ps := []pod{}
	for _, p := range pods.Items {
		ps = append(ps, pod{
			Name:              p.Name,
			Status:            string(p.Status.Phase),
			ContainerStatuses: p.Status.ContainerStatuses,
		})
	}
	return ps

	// // Examples for error handling:
	// // - Use helper functions like e.g. errors.IsNotFound()
	// // - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	// _, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
	// if errors.IsNotFound(err) {
	// 	fmt.Printf("Pod not found\n")
	// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	// } else if err != nil {
	// 	panic(err.Error())
	// } else {
	// 	fmt.Printf("Found pod\n")
	// }
}
