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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func getConfig() (config Config) {
	// Get path to kube config.yaml
	filePtr := flag.String("config", "kube.yml", "path to Kubernetes client config yaml")
	podurl := flag.String("podurl", "", "URL string of voltest pod ingress")
	flag.Parse()
	config.KubeConfigFile = *filePtr
	config.PodUrl = *podurl
	return config
}

func printPVCs(pvcs *v1.PersistentVolumeClaimList) {
	template := "%-32s%-8s%-8s\n"
	fmt.Printf(template, "NAME", "STATUS", "CAPACITY")
	for _, pvc := range pvcs.Items {
		quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]
		fmt.Printf(
			template,
			pvc.Name,
			string(pvc.Status.Phase),
			quant.String())
	}
}

func getContainerCall(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return string(body)
	}
	return ""
}

func main() {
	configVars := getConfig()
	var err error

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

	version, err := c.Discovery().ServerVersion()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Version is %s\n", version)

	// Confirm test pod exists

	namespace := "default"
	pod := "voltest-0"
	// foo := Pod.new()
	//	api := c.CoreV1()
	_, err = c.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	}

	// Confirm test pod is running

	p, err := c.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})

	fmt.Printf("Pod %s is %s\n", pod, p.Status.Phase)

	// Confirm that status page of container is happy
	statusUrl := configVars.PodUrl + "/status"
	fmt.Println(statusUrl)
	resp := getContainerCall(statusUrl)
	if resp == "OK" {
		fmt.Println("Pod Status is Happy")
	}

	// Clear storage data
	// Slight bug here, workaround is explained:
	// the resetfilecheck call should return "1",
	// but my test environment is running an older version of the
	// container. So... for now we'll just run the /textcheck
	// and /bincheck and confirm that we expect "0" for those
	// after running the reset.

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

	resp = getContainerCall(configVars.PodUrl + "/textcheck")
	if resp == "1" {
		fmt.Println("After Reset, textcheck passes as expected")
	} else {
		fmt.Println("Textcheck failed")
	}

	// Confirm binfile

	resp = getContainerCall(configVars.PodUrl + "/bincheck")
	if resp == "1" {
		fmt.Println("After Reset, bincheck passes as expected")
	} else {
		fmt.Println("bincheck failed")
	}

	// Reschedule container
	fmt.Println("Shutting down container")
	//We're not using getContainerCall because an http error here
	// is expected and okay
	sresp, err := http.Get(configVars.PodUrl + "/shutdown")
	if err != nil && sresp != nil {
		fmt.Println("http error okay here")
	}

	// Confirm textfile on rescheduled container
	// This can take a little time, so we'll loop around a sleep

	// We are bypassing this for dev purposes 'cause it takes time
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

	fmt.Println("Confirming container data after restart")
	resp = getContainerCall(configVars.PodUrl + "/textcheck")
	if resp == "1" {
		fmt.Println("After Reset, textcheck passes as expected")
	} else {
		fmt.Println("Textcheck failed")
	}
	resp = getContainerCall(configVars.PodUrl + "/bincheck")
	if resp == "1" {
		fmt.Println("After Reset, bincheck passes as expected")
	} else {
		fmt.Println("bincheck failed")
	}

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

	fmt.Println("Confirming container data after reschedule")
	resp = getContainerCall(configVars.PodUrl + "/textcheck")
	if resp == "1" {
		fmt.Println("After Reset, textcheck passes as expected")
	} else {
		fmt.Println("Textcheck failed")
	}
	resp = getContainerCall(configVars.PodUrl + "/bincheck")
	if resp == "1" {
		fmt.Println("After Reset, bincheck passes as expected")
	} else {
		fmt.Println("bincheck failed")
	}

	// Cleanup post test:
	// reset unschedulable node back to schedulable

	n.Spec.Unschedulable = false
	n, err = c.CoreV1().Nodes().Update(n)

}
