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
	"net"
	"os"
)

type InterfaceDetails struct {
	name      string
	hwaddress string
	addresses []string
}

type PodDetails struct {
	name          string
	netInterfaces []InterfaceDetails
}

func GetPodDetails() (podDetails *PodDetails) {

	podName := os.Getenv("MY_POD_NAME")
	if podName == "" {
		podName = os.Getenv("HOSTNAME")
	}
	podDetails = &PodDetails{name: podName, netInterfaces: make([]InterfaceDetails, 0)}

	l, err := net.Interfaces()

	if err == nil {
		for _, f := range l {
			netInterface := InterfaceDetails{}
			netInterface.name = f.Name
			netInterface.hwaddress = f.HardwareAddr.String()

			addrs, err := f.Addrs()

			if err == nil {
				netInterface.addresses = make([]string, 0)
				for _, addr := range addrs {
					netInterface.addresses = append(netInterface.addresses, addr.String())
				}
			}
			podDetails.netInterfaces = append(podDetails.netInterfaces, netInterface)
		}

	}
	return
}
