package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"net/http"

	"github.com/spf13/viper"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start leanix exporter server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		s := http.NewServeMux()
		s.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {

			b, _ := json.Marshal(leanixExport{
				Namespaces: getNamespaces(clientset),
			})
			fmt.Fprintln(rw, string(b))
		})

		if err := http.ListenAndServe(":8000", s); err != nil {
			panic(err)
		}
	},
}

type pod struct {
	Name              string
	Status            string
	ContainerStatuses []v1.ContainerStatus
}
type namespace struct {
	Name string
	Pods []pod
}
type leanixExport struct {
	Namespaces []namespace
}

func getNamespaces(c *kubernetes.Clientset) []namespace {
	// creates the in-cluster config
	ns, err := c.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	s := []namespace{}
	for _, n := range ns.Items {

		if !isExcluded(n.Name) {
			s = append(s, namespace{
				Name: n.Name,
				Pods: getPods(c, n.Name),
			})
		}
	}

	return s
}

func isExcluded(ns string) bool {
	excludes := viper.GetStringSlice("excludes")
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

func init() {
	RootCmd.AddCommand(runCmd)
}
