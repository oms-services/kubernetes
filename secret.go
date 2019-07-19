package main

import (
	"encoding/json"
	"errors"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// SecretDefinition Required Args as defined in the microservice.yml
type SecretDefinition struct {
	Namespace string        `json:"namespace"`
	Name      string        `json:"name"`
	Secret    corev1.Secret `json:"secret"`
}

// SecretHandler for the core/v1/Secret K8s API
func SecretHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var secretDef SecretDefinition
		json.NewDecoder(r.Body).Decode(&secretDef)
		secretDef.Secret.ObjectMeta.Name = secretDef.Name

		var secret *corev1.Secret
		var err error
		switch r.Method {
		case http.MethodPost:
			secret, err = clientset.CoreV1().Secrets(secretDef.Namespace).Create(&secretDef.Secret)
		case http.MethodPut:
			secret, err = clientset.CoreV1().Secrets(secretDef.Namespace).Update(&secretDef.Secret)
		default:
			panic(errors.New("Unexpected HTTP Method. This should be defined in microservice.yml"))
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(secret)
	}
}
