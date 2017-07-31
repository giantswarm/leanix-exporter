package k8s

import (
	"log"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/batch/v2alpha1"
	"k8s.io/api/core/v1"
	v1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/leanix-exporter/service/exporter/k8s"
)

type metadata struct {
	Labels map[string]string `json:"labels,omitempty"`
}

type Service struct {
	metadata `json:"metadata,omitempty"`

	Name     string            `json:"name,omitempty"`
	Ports    []v1.ServicePort  `json:"ports,omitempty"`
	Type     v1.ServiceType    `json:"type,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
}

type Deployment struct {
	metadata `json:"metadata,omitempty"`

	Name   string                   `json:"name,omitempty"`
	Status v1beta1.DeploymentStatus `json:"status,omitempty"`
}

type Pod struct {
	metadata `json:"metadata,omitempty"`

	Name              string               `json:"name,omitempty"`
	Status            string               `json:"status,omitempty"`
	ContainerStatuses []v1.ContainerStatus `json:"container_statuses,omitempty"`
}

type PodTemplate struct {
	metadata      `json:"metadata,omitempty"`
	Containers    []v1.Container `json:"containers,omitempty"`
	RestartPolicy string         `json:"restart_policy,omitempty"`
	DNSPolicy     string         `json:"dns_policy,omitempty"`
}

type DaemonSet struct {
	metadata    `json:"metadata,omitempty"`
	Status      v1b1.DaemonSetStatus `json:"status,omitempty"`
	PodTemplate PodTemplate          `json:"pod_template,omitempty"`
	Selector    metav1.LabelSelector `json:"selector,omitempty"`
}

type StatefulSet struct {
	metadata    `json:"metadata,omitempty"`
	ServiceName string                    `json:"service_name,omitempty"`
	Replicas    int32                     `json:"replicas,omitempty"`
	PodTemplate PodTemplate               `json:"pod_template,omitempty"`
	Selector    metav1.LabelSelector      `json:"selector,omitempty"`
	Status      v1beta1.StatefulSetStatus `json:"status,omitempty"`
}

type CronJob struct {
	metadata    `json:"metadata,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Schedule    string                 `json:"schedule,omitempty"`
	Status      v2alpha1.CronJobStatus `json:"status,omitempty"`
	Suspend     bool                   `json:"suspend,omitempty"`
	JobTemplate JobTemplate            `json:"job_template,omitempty"`
}

type JobTemplate struct {
	metadata    `json:"metadata,omitempty"`
	PodTemplate PodTemplate `json:"pod_template,omitempty"`
}

type Ingress struct {
	metadata `json:"metadata,omitempty"`
	Backends v1b1.IngressBackend `json:"backends,omitempty"`
	Rules    []v1b1.IngressRule  `json:"rules,omitempty"`
	TLSHosts []v1b1.IngressTLS   `json:"tls_hosts,omitempty"`
	Status   v1b1.IngressStatus  `json:"status,omitempty"`
}

type Namespace struct {
	metadata `json:"metadata,omitempty"`

	Name         string        `json:"name,omitempty"`
	Pods         []Pod         `json:"pods,omitempty"`
	Deployments  []Deployment  `json:"deployments,omitempty"`
	Services     []Service     `json:"services,omitempty"`
	DaemonSets   []DaemonSet   `json:"daemon_sets,omitempty"`
	StatefulSets []StatefulSet `json:"stateful_sets,omitempty"`
	CronJobs     []CronJob     `json:"cron_jobs,omitempty"`
	Ingresses    []Ingress     `json:"ingresses,omitempty"`
}

func FromServiceNamespaces(o []k8s.Namespace) []Namespace {
	ps := []Namespace{}
	for _, p := range o {
		log.Println(p.Labels)
		ps = append(ps, Namespace{
			metadata: metadata{
				Labels: p.Labels,
			},
			Name:         p.Name,
			Pods:         fromServicePods(p.Pods),
			Deployments:  fromServiceDeployments(p.Deployments),
			Services:     fromServiceServices(p.Services),
			DaemonSets:   fromServiceDaemonSets(p.DaemonSet),
			StatefulSets: fromServiceStatefulSets(p.StatefulSets),
			CronJobs:     fromServiceCronJobs(p.CronJobs),
			Ingresses:    fromServiceIngresses(p.Ingresses),
		})
	}

	return ps
}

func fromServicePods(o []k8s.Pod) []Pod {
	ps := []Pod{}
	for _, p := range o {
		ps = append(ps, Pod{
			metadata: metadata{
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
			metadata: metadata{
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
			metadata: metadata{
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

func fromServiceDaemonSets(o []k8s.DaemonSet) []DaemonSet {
	ps := []DaemonSet{}
	for _, p := range o {
		ps = append(ps, DaemonSet{
			metadata: metadata{
				Labels: p.Labels,
			},
			Status:      p.Status,
			Selector:    p.Selector,
			PodTemplate: fromServicePodTemplate(p.PodTemplate),
		})
	}
	return ps
}

func fromServicePodTemplate(o k8s.PodTemplate) PodTemplate {
	return PodTemplate{
		metadata: metadata{
			Labels: o.Labels,
		},
		Containers:    o.Containers,
		DNSPolicy:     o.DNSPolicy,
		RestartPolicy: o.RestartPolicy,
	}
}

func fromServiceStatefulSets(o []k8s.StatefulSet) []StatefulSet {
	ps := []StatefulSet{}
	for _, p := range o {
		ps = append(ps, StatefulSet{
			metadata: metadata{
				Labels: p.Labels,
			},
			ServiceName: p.ServiceName,
			Replicas:    p.Replicas,
			PodTemplate: fromServicePodTemplate(p.PodTemplate),
			Selector:    p.Selector,
			Status:      p.Status,
		})
	}
	return ps
}

func fromServiceCronJobs(o []k8s.CronJob) []CronJob {
	ps := []CronJob{}
	for _, p := range o {
		ps = append(ps, CronJob{
			metadata: metadata{
				Labels: p.Labels,
			},
			Name:        p.Name,
			Schedule:    p.Schedule,
			Status:      p.Status,
			Suspend:     p.Suspend,
			JobTemplate: fromServiceJobTemplate(p.JobTemplate),
		})
	}
	return ps
}

func fromServiceIngresses(o []k8s.Ingress) []Ingress {
	ps := []Ingress{}
	for _, p := range o {
		ps = append(ps, Ingress{
			metadata: metadata{
				Labels: p.Labels,
			},
			Status:   p.Status,
			Backends: p.Backends,
			Rules:    p.Rules,
			TLSHosts: p.TLSHosts,
		})
	}
	return ps
}

func fromServiceJobTemplate(o k8s.JobTemplate) JobTemplate {
	return JobTemplate{
		metadata: metadata{
			Labels: o.Labels,
		},
		PodTemplate: fromServicePodTemplate(o.PodTemplate),
	}
}
