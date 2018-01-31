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
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	//"k8s.io/apimachinery/pkg/api/errors"
)

func main() {

	podIP, hwAddr := getPodNetworkDetails()
	fmt.Printf("The pod IP is: %s\n", podIP)
	fmt.Printf("The pod HW Address is: %s\n", hwAddr)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
}
