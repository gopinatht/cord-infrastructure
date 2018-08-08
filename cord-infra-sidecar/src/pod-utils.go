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
	"net"
	"os"
	"strings"
)

func getPodNetworkDetails() (podIP string, hwAddr string) {

	podIP = os.Getenv("MY_POD_IP")
	hwAddr = ""
	l, err := net.Interfaces()
	found := false

	if err == nil {
		for _, f := range l {

			addrs, err := f.Addrs()

			if err == nil {

				for _, addr := range addrs {

					if strings.Contains(addr.String(), podIP) {
						fmt.Printf("Interface Name: %s\n", f.Name)
						fmt.Printf("Interface Address: %s\n", addr.String())
						fmt.Printf("Interface HW Address: %s\n", f.HardwareAddr.String())
						hwAddr = f.HardwareAddr.String()
						found = true
						break
					}
				}

			}

			if found {
				break
			}

		}
	}
	return
}
