package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"k8s.io/api/core/v1"
)

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
		log.Print(err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
		return string(body)
	}
	return ""
}

func binCheck(podurl string, testname string) TestCheck {
	var test TestCheck
	test.Name = testname
	resp := getContainerCall(podurl + "/bincheck")
	if resp == "1" {
		test.Passed = true
		test.Message = "Bincheck passes as expected"
	} else {
		test.Passed = false
		test.Message = "bincheck failed"
	}
	return test
}

func appendTestCheck(list []TestCheck, test TestCheck) []TestCheck {
	list = append(list, test)
	if test.Passed == true {
	} else {
	}
	return list
}

func textCheck(podurl string, testname string) TestCheck {
	var test TestCheck
	test.Name = testname
	resp := getContainerCall(podurl + "/textcheck")
	if resp == "1" {
		test.Passed = true
		test.Message = "Textcheck passes as expected"
	} else {
		test.Passed = false
		test.Message = "Textcheck failed"
	}
	return test
}

func reportAndOutput(list []TestCheck) int {
	exitcode := 0
	fmt.Println("Test results:")
	fmt.Println("+-------------------------------------------------------+")
	for _, test := range list {
		fmt.Print(test.Name + ": " + test.Message + "\t\t")
		if test.Passed == false {
			exitcode = 1
			fmt.Println("FAILED")
		} else {
			fmt.Println("OK")
		}
	}
	if exitcode == 0 {
		fmt.Println("All tests passed.")
	} else {
		fmt.Println("Check failures.")
	}
	return exitcode
}
