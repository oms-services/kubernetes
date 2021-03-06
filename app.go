package main

import (
	"flag"
	"log"

	"net/http"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/pkg/client"
)

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func getConfigForKube() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	var kubeconfig *string
	home, err := os.UserHomeDir()
	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")

	flag.Parse()

	// use the current context in kubeconfig
	return clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// clientcmd.BuildConfigFromKubeconfigGetter()
	// clientcmd.api.Config
}

func authenticateKube() {
	clientcmd.Load([]byte(os.Getenv("kubeconfig")))
}

func main() {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(os.Getenv("kubeconfig")))

	// config, err := clientcmd.Load([]byte(os.Getenv("kubeconfig")))
	// config, err := getConfigForKube()
	if err != nil {
		panic(err.Error())
	}

	// kubernetes.New
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler())
	mux.Handle("/createJob", CreateJobHandler(clientset))
	mux.Handle("/secret", SecretHandler(clientset))

	log.Print("listening on 8080")
	http.ListenAndServe(":8080", mux)

}
