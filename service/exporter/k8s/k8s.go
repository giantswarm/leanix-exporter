package k8s

import (
	"log"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	v1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
)

type metadata struct {
	Labels map[string]string
}

type Service struct {
	metadata
	Name     string
	Ports    []v1.ServicePort
	Type     v1.ServiceType
	Selector map[string]string
}

type Deployment struct {
	metadata
	Name   string
	Status v1beta1.DeploymentStatus
}

type Pod struct {
	metadata
	Name              string
	Status            string
	ContainerStatuses []v1.ContainerStatus
}

type PodTemplate struct {
	metadata
	Containers    []v1.Container
	RestartPolicy string
	DNSPolicy     string
}

type DaemonSet struct {
	metadata
	Status      v1b1.DaemonSetStatus
	PodTemplate PodTemplate
	Selector    metav1.LabelSelector
}

type Namespace struct {
	metadata
	Name         string
	Pods         []Pod
	Deployments  []Deployment
	Services     []Service
	DaemonSet    []DaemonSet
	StatefulSets []StatefulSet
}

func GetNamespaces(c *kubernetes.Clientset, excludes []string) []Namespace {
	// creates the in-cluster config
	ns, err := c.Namespaces().List(metav1.ListOptions{})
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

			dss, err := getDaemonSets(c, n.Name)
			if err != nil {
				log.Println(err)
			}

			sss, err := getStatefulSets(c, n.Name)
			if err != nil {
				log.Println(err)
			}

			s = append(s, Namespace{
				metadata: metadata{
					Labels: n.Labels,
				},
				Name:         n.Name,
				Pods:         getPods(c, n.Name),
				Deployments:  depls,
				Services:     svcs,
				DaemonSet:    dss,
				StatefulSets: sss,
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
			metadata: metadata{
				Labels: d.Labels,
			},
		})
	}

	return s, nil
}

func getPods(c *kubernetes.Clientset, ns string) []Pod {
	// creates the in-cluster config
	pods, err := c.Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	ps := []Pod{}
	for _, p := range pods.Items {
		ps = append(ps, Pod{
			Name:              p.Name,
			Status:            string(p.Status.Phase),
			ContainerStatuses: p.Status.ContainerStatuses,
			metadata: metadata{
				Labels: p.Labels,
			},
		})
	}
	return ps
}

func getServices(c *kubernetes.Clientset, ns string) ([]Service, error) {
	services, err := c.Services(ns).List(metav1.ListOptions{})
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
			metadata: metadata{
				Labels: s.Labels,
			},
		})
	}
	return ss, nil
}

func getDaemonSets(c *kubernetes.Clientset, ns string) ([]DaemonSet, error) {
	dss, err := c.DaemonSets(ns).List(metav1.ListOptions{})

	if err != nil {
		return []DaemonSet{}, microerror.Mask(err)
	}

	ds := []DaemonSet{}
	for _, s := range dss.Items {

		ds = append(ds, DaemonSet{
			metadata: metadata{
				Labels: s.Labels,
			},
			Status:      s.Status,
			Selector:    *s.Spec.Selector,
			PodTemplate: fromPodTemplateSpec(s.Spec.Template),
		})
	}
	return ds, nil
}

func fromPodTemplateSpec(pst v1.PodTemplateSpec) PodTemplate {
	return PodTemplate{
		metadata: metadata{
			Labels: pst.Labels,
		},
		Containers:    pst.Spec.Containers,
		RestartPolicy: string(pst.Spec.RestartPolicy),
		DNSPolicy:     string(pst.Spec.DNSPolicy),
	}
}

type StatefulSet struct {
	metadata
	ServiceName string
	Replicas    int32
	PodTemplate PodTemplate
	Selector    metav1.LabelSelector
	Status      v1beta1.StatefulSetStatus
}

func getStatefulSets(c *kubernetes.Clientset, ns string) ([]StatefulSet, error) {

	ss := []StatefulSet{}
	statefulSets, err := c.AppsV1beta1().StatefulSets(ns).List(metav1.ListOptions{})

	if err != nil {
		return ss, microerror.Mask(err)
	}

	for _, s := range statefulSets.Items {

		ss = append(ss, StatefulSet{
			metadata: metadata{
				Labels: s.Labels,
			},
			ServiceName: s.Spec.ServiceName,
			Replicas:    *s.Spec.Replicas,
			PodTemplate: fromPodTemplateSpec(s.Spec.Template),
			Selector:    *s.Spec.Selector,
			Status:      s.Status,
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
