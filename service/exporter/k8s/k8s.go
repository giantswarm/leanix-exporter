package k8s

import (
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/batch/v2alpha1"
	"k8s.io/api/core/v1"
	v1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/microerror"
	micrologger "github.com/giantswarm/microkit/logger"
)

type metadata struct {
	Labels map[string]string
}

// Service is the kubernetes Service data
type Service struct {
	metadata
	Name     string
	Ports    []v1.ServicePort
	Type     v1.ServiceType
	Selector map[string]string
}

// Deployment is the kubernetes Deployment data
type Deployment struct {
	metadata
	Name   string
	Status v1beta1.DeploymentStatus
}

// Pod is the kubernetes Pod data
type Pod struct {
	metadata
	Name              string
	Status            string
	ContainerStatuses []v1.ContainerStatus
}

// PodTemplate is the kubernetes PodTemplate data
type PodTemplate struct {
	metadata
	Containers    []v1.Container
	RestartPolicy string
	DNSPolicy     string
}

// DaemonSet is the kubernetes DaemonSet data
type DaemonSet struct {
	metadata
	Status      v1b1.DaemonSetStatus
	PodTemplate PodTemplate
	Selector    metav1.LabelSelector
}

// StatefulSet is the kubernetes StatefulSet data
type StatefulSet struct {
	metadata
	ServiceName string
	Replicas    int32
	PodTemplate PodTemplate
	Selector    metav1.LabelSelector
	Status      v1beta1.StatefulSetStatus
}

// CronJob is the kubernetes CronJob data
type CronJob struct {
	metadata
	Name        string
	Schedule    string
	Status      v2alpha1.CronJobStatus
	Suspend     bool
	JobTemplate JobTemplate
}

// JobTemplate is the kubernetes JobTemplate data
type JobTemplate struct {
	metadata
	PodTemplate PodTemplate
}

// Ingress is the kubernetes Ingress data
type Ingress struct {
	metadata
	Backends *v1b1.IngressBackend
	Rules    []v1b1.IngressRule
	TLSHosts []v1b1.IngressTLS
	Status   v1b1.IngressStatus
}

// NetworkPolicy is the kubernetes NetworkPolicy data
type NetworkPolicy struct {
	metadata
	Name         string
	IngressRules []v1b1.NetworkPolicyIngressRule
	Selector     metav1.LabelSelector
}

// Namespace is the kubernetes Namespace data
type Namespace struct {
	metadata
	Name            string
	Pods            []Pod
	Deployments     []Deployment
	Services        []Service
	DaemonSet       []DaemonSet
	StatefulSets    []StatefulSet
	CronJobs        []CronJob
	Ingresses       []Ingress
	NetworkPolicies []NetworkPolicy
}

// GetNamespaces returns the kubernetes Namespace related aggregated data
func GetNamespaces(c *kubernetes.Clientset, excludes []string, log micrologger.Logger) ([]Namespace, error) {
	ns, err := c.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	s := []Namespace{}
	for _, n := range ns.Items {
		if !isExcluded(excludes, n.Name) {
			depls, err := getDeployments(c, n.Name)
			if err != nil {
				log.Log("warning", err)
			}

			pods, err := getPods(c, n.Name)
			if err != nil {
				log.Log("warning", err)
			}

			svcs, err := getServices(c, n.Name)
			if err != nil {
				log.Log("warning", err)
			}

			dss, err := getDaemonSets(c, n.Name)
			if err != nil {
				log.Log("warning", err)
			}

			sss, err := getStatefulSets(c, n.Name)
			if err != nil {
				log.Log("warning", err)
			}
			cjs, err := getCronJobs(c, n.Name)
			if err != nil {
				log.Log("warning", err)
			}

			is, err := getIngresses(c, n.Name)
			if err != nil {
				log.Log("warning", err)
			}
			nps, err := getNetworkPolicies(c, n.Name)
			if err != nil {
				log.Log("warning", err)
			}

			s = append(s, Namespace{
				metadata: metadata{
					Labels: n.Labels,
				},
				Name:            n.Name,
				Pods:            pods,
				Deployments:     depls,
				Services:        svcs,
				DaemonSet:       dss,
				StatefulSets:    sss,
				CronJobs:        cjs,
				Ingresses:       is,
				NetworkPolicies: nps,
			})
		}
	}

	return s, nil
}

func getDeployments(c *kubernetes.Clientset, ns string) ([]Deployment, error) {
	depls, err := c.AppsV1beta1().Deployments(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
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

func getPods(c *kubernetes.Clientset, ns string) ([]Pod, error) {
	pods, err := c.Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
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
	return ps, nil
}

func getServices(c *kubernetes.Clientset, ns string) ([]Service, error) {
	services, err := c.Services(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
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
		return nil, microerror.Mask(err)
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

func getStatefulSets(c *kubernetes.Clientset, ns string) ([]StatefulSet, error) {
	statefulSets, err := c.AppsV1beta1().StatefulSets(ns).List(metav1.ListOptions{})

	if err != nil {
		return nil, microerror.Mask(err)
	}
	ss := []StatefulSet{}

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

func getCronJobs(c *kubernetes.Clientset, ns string) ([]CronJob, error) {
	cronjobs, err := c.BatchV2alpha1().CronJobs(ns).List(metav1.ListOptions{})

	if err != nil {
		return nil, microerror.Mask(err)
	}
	cj := []CronJob{}

	for _, s := range cronjobs.Items {
		cj = append(cj, CronJob{
			metadata: metadata{
				Labels: s.Labels,
			},
			Name:        s.GetName(),
			Schedule:    s.Spec.Schedule,
			Status:      s.Status,
			Suspend:     *s.Spec.Suspend,
			JobTemplate: fromJobTemplateSpec(s.Spec.JobTemplate),
		})
	}
	return cj, nil
}

func getIngresses(c *kubernetes.Clientset, ns string) ([]Ingress, error) {
	ingresses, err := c.ExtensionsV1beta1().Ingresses(ns).List(metav1.ListOptions{})

	if err != nil {
		return nil, microerror.Mask(err)
	}
	is := []Ingress{}
	for _, i := range ingresses.Items {
		is = append(is, Ingress{
			metadata: metadata{
				Labels: i.Labels,
			},
			Status:   i.Status,
			Backends: i.Spec.Backend,
			Rules:    i.Spec.Rules,
			TLSHosts: i.Spec.TLS,
		})
	}

	return is, nil
}

func getNetworkPolicies(c *kubernetes.Clientset, ns string) ([]NetworkPolicy, error) {
	npl := &v1b1.NetworkPolicyList{}
	err := c.ExtensionsV1beta1().RESTClient().Get().
		Namespace(ns).
		Resource("networkpolicies").
		Do().
		Into(npl)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	nps := []NetworkPolicy{}
	for _, n := range npl.Items {
		nps = append(nps, NetworkPolicy{
			metadata: metadata{
				Labels: n.Labels,
			},
			Name:         n.Name,
			IngressRules: n.Spec.Ingress,
			Selector:     n.Spec.PodSelector,
		})
	}

	return nps, nil
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

func fromJobTemplateSpec(pst v2alpha1.JobTemplateSpec) JobTemplate {
	return JobTemplate{
		metadata: metadata{
			Labels: pst.Labels,
		},
		PodTemplate: fromPodTemplateSpec(pst.Spec.Template),
	}
}

func isExcluded(excludes []string, ns string) bool {
	for _, e := range excludes {
		if e == ns {
			return true
		}
	}

	return false
}
