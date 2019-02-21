## Introduction

This is the Kubernetes-specific certification system needed to test storage for Docker support. The test is designed to work against any valid storage driver, no matter which API. Follow instructions below to get started. If you need help or have feedback, contact me at keith@docker.com or through GitHub.

The tests revolve around a container that runs most of the important bits. The container is deployed in a Kubernetes Service using the NodePort network access method, so once you're deployed, you can hit any node in the Kube cluster on that port via HTTP and have access to the container's native API.

## Setup

Currently, this only works with in-tree drivers (Early access versions of Docker Enterprise have CSI and FlexColume as well).

Install Docker EE and configure your storage integration. Define a StorageClass based on your storage driver. The latest Docker EE ships with Kubernetes 1.11, and has from Kube beta-level CSI support, as well as GA FlexVolume. If you want to test a FlexVolume or CSI plugin, please contact us at the above email.

To get going, edit the kubernetes/testapp.yaml file. The only config change needed is to provide your storageClass.

Then, run the following:

```
kubectl apply -f kubernetes/testapp.yaml
```

Once the deployment is running, you can find the publish port by running:

```
kubectl describe service voltestservice
```

Now you can kick the container around and exercise storage using the container's lightweight API. That's documented in the containerdocs.md file.

The actual test program is in early development - it works, and you can run tests, but formatting is not a current concern at the moment. Still, feel free to build and run it! It's currently one self-contained Golang file, so you can build it with:

```
go build voltestkube.go voltestkubelib.go
```

To run the tests, you need two things:

# Path to a working Kubernetes client config file (kube.yml, etc)
# Url to your pod's NodePort

Invoke the program like this:

```
./voltestkube --config="/path/to/kube.yml" --podurl="http://10.2.2.74:32779"
```

The config path does NOT expand, so don't use `~/` shortcuts - full path only for now.

Test output will look like this if everything's working:

```
Pod voltest-0 is Running
http://3.17.176.248:34511/status
Reset test data for clean run
Shutting down container
Waiting for container restart - we wait up to 10 minutes
Should be pulling status from http://3.17.176.248:34511/status
.Get http://3.17.176.248:34511/status: dial tcp 3.17.176.248:34511: connect: connection refused
.Get http://3.17.176.248:34511/status: dial tcp 3.17.176.248:34511: connect: connection refused
.Get http://3.17.176.248:34511/status: dial tcp 3.17.176.248:34511: connect: connection refused
.Get http://3.17.176.248:34511/status: dial tcp 3.17.176.248:34511: connect: connection refused
Container restarted successfully, moving on
Pod node voltest-0 is ip-172-31-1-136.us-east-2.compute.internal
Pod was running on ip-172-31-1-136.us-east-2.compute.internal
Shutting down container for forced reschedule
http error okay here
Waiting for container rechedule - we wait up to 10 minutes
.....................Container rescheduled successfully, moving on
Pod is now running on ip-172-31-14-222.us-east-2.compute.internal
Going into cleanup...
Cleaning up taint on ip-172-31-1-136.us-east-2.compute.internal
Test results:
+-------------------------------------------------------+
Kubernetes Version: v1.11.5-docker-1-golang-1106		OK
Test Pod Existence: Found pod voltest-0 in namespace default		OK
Confirm Running Pod: Pod running		OK
Initial Textfile Content Confirmation: Textcheck passes as expected		OK
Initial Binary Content Confirmation: Bincheck passes as expected		OK
Post-restart Textfile Content Confirmation: Textcheck passes as expected		OK
Post-restart Binaryfile Content Confirmation: Bincheck passes as expected		OK
Rescheduled Textfile Content Confirmation: Textcheck passes as expected		OK
Rescheduled Binaryfile Content Confirmation: Bincheck passes as expected		OK
All tests passed.
```

Critical notices are "Passes as expected". Refer to the voltestkube.go code for more info.


## Cheat sheet commands:

Watch container status/restart loop:
kubectl get pod -l app=voltest -w

Upgrade container in running deployment:

kubectl patch statefulset voltest --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/image", "value":"khudgins/volcheck:latest"}]'

kubectl describe service voltest
kubectl describe statefulset voltest

./voltestkube --config="/full/path/to/kube.yml"
