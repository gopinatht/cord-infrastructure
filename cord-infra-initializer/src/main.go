// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghodss/yaml"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	defaultAnnotation      = "initializer.onf.io/cord-infra"
	defaultInitializerName = "cord-infra.initializer.onf.io"
	defaultConfigmap       = "cord-infra-initializer"
	defaultNamespace       = "utility"
)

var (
	annotation        string
	configmap         string
	initializerName   string
	namespace         string
	requireAnnotation bool
)

type config struct {
	Containers []corev1.Container
	Volumes    []corev1.Volume
}

func main() {
	flag.StringVar(&annotation, "annotation", defaultAnnotation, "The annotation to trigger initialization")
	flag.StringVar(&configmap, "configmap", defaultConfigmap, "The cord-infra initializer configuration configmap")
	flag.StringVar(&initializerName, "initializer-name", defaultInitializerName, "The initializer name")
	flag.StringVar(&namespace, "namespace", defaultNamespace, "The configuration namespace")
	flag.BoolVar(&requireAnnotation, "require-annotation", false, "Require annotation for initialization")
	flag.Parse()

	log.Println("Starting the CORD Infrastructure initializer...")
	log.Printf("Initializer name set to: %s", initializerName)

	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Load the CORD Infrastructure Initializer configuration from a Kubernetes ConfigMap.
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(configmap, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	c, err := configmapToConfig(cm)
	if err != nil {
		log.Fatal(err)
	}

	// Watch uninitialized Deployments, Replicasets in all namespaces.
	restClient := clientset.AppsV1beta2().RESTClient()
	deploymentStop := make(chan struct{})
	replicasetStop := make(chan struct{})
	watchDeployment(c, clientset, restClient, deploymentStop)
	watchReplicaSet(c, clientset, restClient, replicasetStop)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	close(deploymentStop)
	close(replicasetStop)
}

func configmapToConfig(configmap *corev1.ConfigMap) (*config, error) {
	var c config
	err := yaml.Unmarshal([]byte(configmap.Data["config"]), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
