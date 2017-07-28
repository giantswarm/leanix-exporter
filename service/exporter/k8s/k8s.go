package k8s

import (
	"log"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
)

type obj struct {
	Labels map[string]string
}

type Service struct {
	obj
	Name     string
	Ports    []v1.ServicePort
	Type     v1.ServiceType
	Selector map[string]string
}

type Deployment struct {
	obj
	Name   string
	Status v1beta1.DeploymentStatus
}

type Pod struct {
	obj
	Name              string
	Status            string
	ContainerStatuses []v1.ContainerStatus
}
type Namespace struct {
	obj
	Name        string
	Pods        []Pod
	Deployments []Deployment
	Services    []Service
}

func GetNamespaces(c *kubernetes.Clientset, excludes []string) []Namespace {
	// creates the in-cluster config
	ns, err := c.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	s := []Namespace{}
	for _, n := range ns.Items {
		if !isExcluded(excludes, n.Name) {
			depls, err := getDeployments(c, n.Name)
			if err != nil {
				log.Println(err)
			}

			svcs, err := getServices(c, n.Name)
			if err != nil {
				log.Println(err)
			}
			s = append(s, Namespace{
				Name:        n.Name,
				Pods:        getPods(c, n.Name),
				Deployments: depls,
				Services:    svcs,
				obj: obj{
					Labels: n.Labels,
				},
			})
		}
	}

	return s
}

func getDeployments(c *kubernetes.Clientset, ns string) ([]Deployment, error) {
	depls, err := c.AppsV1beta1().Deployments(ns).List(metav1.ListOptions{})
	if err != nil {
		return []Deployment{}, microerror.Mask(err)
	}

	s := []Deployment{}
	for _, d := range depls.Items {
		s = append(s, Deployment{
			Name:   d.GetName(),
			Status: d.Status,
			obj: obj{
				Labels: d.Labels,
			},
		})
	}

	return s, nil
}

func getPods(c *kubernetes.Clientset, ns string) []Pod {
	// creates the in-cluster config
	pods, err := c.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	ps := []Pod{}
	for _, p := range pods.Items {

		ps = append(ps, Pod{
			Name:              p.Name,
			Status:            string(p.Status.Phase),
			ContainerStatuses: p.Status.ContainerStatuses,
			obj: obj{
				Labels: p.Labels,
			},
		})
	}
	return ps
}

func getServices(c *kubernetes.Clientset, ns string) ([]Service, error) {
	services, err := c.CoreV1().Services(ns).List(metav1.ListOptions{})
	if err != nil {
		return []Service{}, microerror.Mask(err)
	}

	ss := []Service{}
	for _, s := range services.Items {
		ss = append(ss, Service{
			Name:     s.Name,
			Ports:    s.Spec.Ports,
			Type:     s.Spec.Type,
			Selector: s.Spec.Selector,
			obj: obj{
				Labels: s.Labels,
			},
		})
	}
	return ss, nil
}

func isExcluded(excludes []string, ns string) bool {
	for _, e := range excludes {
		if e == ns {
			return true
		}
	}

	return false
}
