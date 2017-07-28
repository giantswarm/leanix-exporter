package k8s

import (
	"log"

	"github.com/giantswarm/leanix-exporter/service/exporter/k8s"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
)

type obj struct {
	Labels map[string]string `json:"labels,omitempty"`
}

type Service struct {
	obj

	Name     string            `json:"name,omitempty"`
	Ports    []v1.ServicePort  `json:"ports,omitempty"`
	Type     v1.ServiceType    `json:"type,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
}

type Deployment struct {
	obj

	Name   string                   `json:"name,omitempty"`
	Status v1beta1.DeploymentStatus `json:"status,omitempty"`
}

type Pod struct {
	obj

	Name              string               `json:"name,omitempty"`
	Status            string               `json:"status,omitempty"`
	ContainerStatuses []v1.ContainerStatus `json:"container_statuses,omitempty"`
}
type Namespace struct {
	obj

	Name        string       `json:"name,omitempty"`
	Pods        []Pod        `json:"pods,omitempty"`
	Deployments []Deployment `json:"deployments,omitempty"`
	Services    []Service    `json:"services,omitempty"`
}

func FromServiceNamespaces(o []k8s.Namespace) []Namespace {
	ps := []Namespace{}
	for _, p := range o {
		log.Println(p.Labels)
		ps = append(ps, Namespace{
			obj: obj{
				Labels: p.Labels,
			},
			Name:        p.Name,
			Pods:        fromServicePods(p.Pods),
			Deployments: fromServiceDeployments(p.Deployments),
			Services:    fromServiceServices(p.Services),
		})
	}

	return ps
}

func fromServicePods(o []k8s.Pod) []Pod {
	ps := []Pod{}
	for _, p := range o {
		ps = append(ps, Pod{
			obj: obj{
				Labels: p.Labels,
			},
			Name:              p.Name,
			Status:            p.Status,
			ContainerStatuses: p.ContainerStatuses,
		})
	}
	return ps
}

func fromServiceDeployments(o []k8s.Deployment) []Deployment {
	ps := []Deployment{}
	for _, p := range o {
		ps = append(ps, Deployment{
			obj: obj{
				Labels: p.Labels,
			},
			Name:   p.Name,
			Status: p.Status,
		})
	}
	return ps
}
func fromServiceServices(o []k8s.Service) []Service {
	ps := []Service{}
	for _, p := range o {
		ps = append(ps, Service{
			obj: obj{
				Labels: p.Labels,
			},
			Name:     p.Name,
			Ports:    p.Ports,
			Selector: p.Selector,
			Type:     p.Type,
		})
	}
	return ps
}
