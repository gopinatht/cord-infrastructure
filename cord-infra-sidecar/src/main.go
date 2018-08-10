/*
Copyright 2017 Gopinath Taget.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	//"k8s.io/apimachinery/pkg/api/errors"
)

const (
	defaultKafkaBroker = "cord-kafka"
	defaultKafkaTopic  = "dp-pod-details"
)

var (
	kafkaBroker string
	kafkaTopic  string
)

func main() {

	flag.StringVar(&kafkaBroker, "kafka-broker", defaultKafkaBroker, "The kafka broker to use")
	flag.StringVar(&kafkaTopic, "kafka-topic", defaultKafkaTopic, "The kafka topic to use")

	podDetails := GetPodDetails()
	fmt.Printf("The pod details are: %v\n", podDetails)
	podDetailsJSON, err := json.Marshal(podDetails)

	if err == nil {
		fmt.Printf("Sending struct to kafka\n")
		ProduceKafkaMessage(kafkaBroker, kafkaTopic, string(podDetailsJSON))
		fmt.Printf("After Sending struct to kafka\n")

	} else {
		fmt.Printf("\nCould not marshal Struct to Json: %v\n", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
}
