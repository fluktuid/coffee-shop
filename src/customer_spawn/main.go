package main

import (
	"context"
	"io/ioutil"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"

	"github.com/fluktuid/coffee-shop/src/dto"
	"github.com/fluktuid/coffee-shop/src/metrics"
	"github.com/fluktuid/coffee-shop/src/util"
)

func main() {
	var config dto.Config
	util.LoadEnv(&config)
	go metrics.StartMetricsBlock()

	for range time.Tick(time.Second) {
		count := rand.Intn(config.SpawnMax)
		log.Infof("spawn %d jobs", count)
		go createJobs(count, config.CustomerImg)
	}
}

func createJobs(count int, img string) {
	ns := namespace()
	clientSet := connect()
	job(clientSet, img, ns)
}

func connect() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic("Failed to create k8s config")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic("Failed to create k8s clientset")
	}

	return clientset
}

func namespace() string {
	ns, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Panic("Failed to get namespace")
	}
	return string(ns)
}

func job(clientset *kubernetes.Clientset, image, namespace string) {
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "customer-",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "customer-",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "customer",
							Image: image,
							// todo: add envFrom
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
	jobs := clientset.BatchV1().Jobs(namespace)

	_, err := jobs.Create(context.Background(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s job.")
	}

	//print job details
	log.Println("Created K8s job successfully")
}
