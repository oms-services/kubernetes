package main

import (
	"encoding/json"
	"net/http"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

// JobDefinition as defined in the microservice.yml
type JobDefinition struct {
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Spec      batchv1.JobSpec `json:"spec"`
}

// CreateJobHandler creates a K8s Job
func CreateJobHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jobDef JobDefinition
		json.NewDecoder(r.Body).Decode(&jobDef)

		k8Job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name: jobDef.Name,
			},
			Spec: jobDef.Spec,
		}

		v1job, err := clientset.BatchV1().Jobs(jobDef.Namespace).Create(k8Job)
		if err != nil {
			// job creation succeeding has nothing to do with pod success
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(v1job)
	}
}
