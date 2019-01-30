/*
Copyright 2018 Docker, Inc

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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"net/http"
)

type Config struct {
	KubeConfigFile string
	PodUrl         string
}

type TestCheck struct {
	Name    string
	Passed  bool
	Message string
}

var testList []TestCheck
var exitStatus int

type patchBoolValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value bool   `json:"value"`
}

func appendTestCheck(list []TestCheck, test TestCheck) {
	list = append(list, test)
	fmt.Print(test.Name, ": ", test.Message)
	if test.Passed == true {
		fmt.Println("\tOK")
	} else {
		fmt.Println("\tFAIL")
	}
}

func main() {
	exitStatus = 0
	configVars := getConfig()
	var err error
	var test TestCheck

	config, err := clientcmd.BuildConfigFromFlags("", configVars.KubeConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Test coverage:

	// Confirm Kube version
	test.Name = "Kubernetes Version"
	version, err := c.Discovery().ServerVersion()
	if err != nil {
		test.Passed = false
		log.Println(err)
	} else {
		test.Passed = true
		test.Message = version.String()
	}
	appendTestCheck(testList, test)

	// Confirm test pod exists

	namespace := "default"
	pod := "voltest-0"
	test.Name = "Test Pod Existence"
	_, err = c.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		test.Passed = false
		test.Message = "Pod " + pod + " in namespace " + namespace + " not found"
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		test.Passed = false
		test.Message = "Error getting pod " + pod + "in namespace " +
			namespace + ": " + statusError.ErrStatus.Message
	} else if err != nil {
		test.Passed = false
		test.Message = err.Error()
	} else {
		test.Passed = true
		test.Message = "Found pod " + pod + " in namespace " + namespace
	}
	appendTestCheck(testList, test)

	// Confirm test pod is running

	test.Name = "Confirm Running Pod"
	p, err := c.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})

	fmt.Printf("Pod %s is %s\n", pod, p.Status.Phase)

	// Confirm that status page of container is happy
	statusUrl := configVars.PodUrl + "/status"
	fmt.Println(statusUrl)
	resp := getContainerCall(statusUrl)
	if resp == "OK" {
		test.Passed = true
		test.Message = "Pod running"
	} else {
		test.Passed = false
		test.Message = "Pod not running"
	}
	appendTestCheck(testList, test)

	// Clear storage data
	// Slight bug here, workaround is explained:
	// the resetfilecheck call should return "1",
	// but my test environment is running an older version of the
	// container. So... for now we'll just run the /textcheck
	// and /bincheck and confirm that we expect "0" for those
	// after running the reset.

	// Refactor note: This isn't a test - might should be, but we're ignoring that
	// until I get the test container pushed into DockerHub

	getContainerCall(configVars.PodUrl + "/resetfilecheck")
	resp = getContainerCall(configVars.PodUrl + "/textcheck")
	if resp == "0" {
		fmt.Println("After Reset, textcheck fails as expected")
	} else {
		fmt.Println("Something wrong with environment reset, check your environment")
	}
	resp = getContainerCall(configVars.PodUrl + "/bincheck")
	if resp == "0" {
		fmt.Println("After Reset, bincheck fails as expected")
	} else {
		fmt.Println("Something wrong with environment reset, check your environment")
	}

	// Initialize storage data

	getContainerCall(configVars.PodUrl + "/runfilecheck")

	// Confirm textfile

	test = textCheck(configVars.PodUrl, "Initial Textfile Content Confirmation")
	appendTestCheck(testList, test)
	// Confirm binfile

	test = binCheck(configVars.PodUrl, "Initial Binary Content Confirmation")
	appendTestCheck(testList, test)

	// Reschedule container

	//We're not using getContainerCall because an http error here
	// is expected and okay
	// refactor note: Should check to see if we can run getContainerCall now
	fmt.Println("Shutting down container")
	sresp, err := http.Get(configVars.PodUrl + "/shutdown")
	if err != nil && sresp != nil {
		fmt.Println("http error okay here")
	}

	// Confirm textfile on rescheduled container
	// This can take a little time, so we'll loop around a sleep
	//
	fmt.Println("Waiting for container restart - we wait up to 10 minutes")
	fmt.Println("Should be pulling status from " + configVars.PodUrl + "/status")
	for i := 0; i < 60; i++ {
		time.Sleep(10 * time.Second)
		hresp, err := http.Get(configVars.PodUrl + "/status")
		if err != nil {
			fmt.Print(".")
			fmt.Println(err.Error())
		} else {
			body, err := ioutil.ReadAll(hresp.Body)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				if string(body) == "OK" {
					fmt.Println("Container restarted successfully, moving on")
					break
				}
			}
		}
	}

	// confirm binfile on rescheduled container
	test = textCheck(configVars.PodUrl, "Post-restart Textfile Content Confirmation")
	appendTestCheck(testList, test)

	test = binCheck(configVars.PodUrl, "Post-restart Binaryfile Content Confirmation")
	appendTestCheck(testList, test)

	// Force failover test onto a different node.
	// First, let's get the node name:
	// Then, set the node unschedulable
	// then, kill the container

	p, err = c.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	fmt.Printf("Pod node %s is %s\n", pod, p.Spec.NodeName)
	n, err := c.CoreV1().Nodes().Get(p.Spec.NodeName, metav1.GetOptions{})
	n.Spec.Unschedulable = true
	n, err = c.CoreV1().Nodes().Update(n)
	fmt.Println("Pod was running on " + p.Spec.NodeName)
	fmt.Println("Shutting down container for forced reschedule")
	// We're not using getContainerCall because an http error here
	// is expected and okay
	fmt.Println("http error okay here")
	gracePeriodSeconds := int64(0)
	err = c.CoreV1().Pods(namespace).Delete(p.Name, &metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds})

	// Confirm textfile on rescheduled container
	// This can take a little time, so we'll loop around a sleep

	fmt.Println("Waiting for container rechedule - we wait up to 10 minutes")
	for i := 0; i < 60; i++ {
		time.Sleep(10 * time.Second)
		hresp, err := http.Get(configVars.PodUrl + "/status")
		if err != nil {
			fmt.Print(".")
		} else {
			body, err := ioutil.ReadAll(hresp.Body)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				if string(body) == "OK" {
					fmt.Println("Container rescheduled successfully, moving on")
					break
				}
			}
		}
	}
	p, err = c.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	fmt.Println("Pod is now running on " + p.Spec.NodeName)
	// confirm binfile on rescheduled container

	test = textCheck(configVars.PodUrl, "Rescheduled Textfile Content Confirmation")
	appendTestCheck(testList, test)

	test = binCheck(configVars.PodUrl, "Rescheduled Binaryfile Content Confirmation")
	appendTestCheck(testList, test)

	// Cleanup post test:
	// reset unschedulable node back to schedulable
	// They had to make things this complicated, huh?
	fmt.Println("Going into cleanup...")
	//n.Spec.Unschedulable = false
	fmt.Println("Cleaning up taint on " + n.Name)
	patchVal := []patchBoolValue{{
		Op:    "replace",
		Path:  "/spec/unschedulable",
		Value: false,
	}}
	patchValBytes, _ := json.Marshal(patchVal)
	//'{"spec": {"unschedulable": false}}'
	n, err = c.CoreV1().Nodes().Patch(n.GetName(), types.JSONPatchType, patchValBytes)

	os.Exit(exitStatus)
}
