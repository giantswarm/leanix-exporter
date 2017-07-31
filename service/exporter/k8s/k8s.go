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

type StatefulSet struct {
	metadata
	ServiceName string
	Replicas    int32
	PodTemplate PodTemplate
	Selector    metav1.LabelSelector
	Status      v1beta1.StatefulSetStatus
}

type CronJob struct {
	metadata
	Name        string
	Schedule    string
	Status      v2alpha1.CronJobStatus
	Suspend     bool
	JobTemplate JobTemplate
}

type JobTemplate struct {
	metadata
	PodTemplate PodTemplate
}

type Ingress struct {
	metadata
	Backends v1b1.IngressBackend
	Rules    []v1b1.IngressRule
	TLSHosts []v1b1.IngressTLS
	Status   v1b1.IngressStatus
}

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

func GetNamespaces(c *kubernetes.Clientset, excludes []string, log micrologger.Logger) []Namespace {
	// creates the in-cluster config
	ns, err := c.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		//TODO use other than panic
		panic(err.Error())
	}
	s := []Namespace{}
	for _, n := range ns.Items {
		if !isExcluded(excludes, n.Name) {
			depls, err := getDeployments(c, n.Name)
			if err != nil {
				log.Log("debug", err)
			}

			svcs, err := getServices(c, n.Name)
			if err != nil {
				log.Log("debug", err)
			}

			dss, err := getDaemonSets(c, n.Name)
			if err != nil {
				log.Log("debug", err)
			}

			sss, err := getStatefulSets(c, n.Name)
			if err != nil {
				log.Log("debug", err)
			}
			cjs, err := getCronJobs(c, n.Name)
			if err != nil {
				log.Log("debug", err)
			}

			is, err := getIngresses(c, n.Name)
			if err != nil {
				log.Log("debug", err)
			}
			nps, err := GetNetworkPolicies(c, n.Name)
			if err != nil {
				log.Log("debug", err)
			}

			s = append(s, Namespace{
				metadata: metadata{
					Labels: n.Labels,
				},
				Name:            n.Name,
				Pods:            getPods(c, n.Name),
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
		//TODO use other than panic
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

func getCronJobs(c *kubernetes.Clientset, ns string) ([]CronJob, error) {

	cj := []CronJob{}
	cronjobs, err := c.BatchV2alpha1().CronJobs(ns).List(metav1.ListOptions{})

	if err != nil {
		return cj, microerror.Mask(err)
	}

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
	is := []Ingress{}
	ingresses, err := c.ExtensionsV1beta1().Ingresses(ns).List(metav1.ListOptions{})

	if err != nil {
		return is, microerror.Mask(err)
	}
	for _, i := range ingresses.Items {
		is = append(is, Ingress{
			metadata: metadata{
				Labels: i.Labels,
			},
			Status:   i.Status,
			Backends: MustBackend(i.Spec.Backend),
			Rules:    i.Spec.Rules,
			TLSHosts: i.Spec.TLS,
		})
	}

	return is, nil
}

type NetworkPolicy struct {
	metadata
	Name         string
	IngressRules []v1b1.NetworkPolicyIngressRule
	Selector     metav1.LabelSelector
}

func GetNetworkPolicies(c *kubernetes.Clientset, ns string) ([]NetworkPolicy, error) {
	nps := []NetworkPolicy{}
	npl := &v1b1.NetworkPolicyList{}
	err := c.ExtensionsV1beta1().RESTClient().Get().
		Namespace(ns).
		Resource("networkpolicies").
		Do().
		Into(npl)
	if err != nil {
		return []NetworkPolicy{}, microerror.Mask(err)
	}

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

func MustBackend(b *v1b1.IngressBackend) v1b1.IngressBackend {
	if b == nil {
		return v1b1.IngressBackend{}
	}

	return *b
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
