## Introduction

This is the Kubernetes-specific stuff needed to test storage for Docker support. This is currently under development, is ROUGH, but if you want to kick the tires on it, feel free. Just provide feedback to keith@docker.com or through GitHub.

The tests revolve around a container that runs most of the important bits. The container is deployed in a Kubernetes Service using the NodePort network access method, so once you're deployed, you can hit any node in the Kube cluster on that port via HTTP and have access to the container's native API.

## Setup

Currently, this only works with in-tree drivers. You'll need to get Docker EE set up with some type of persistent storage. Before deploying the test framework bits, you'll also need a Kubernetes StorageClass defined in your environment.

To get going, edit the kubernetes/testapp.yaml file. The only config change needed is to provide your storageClass.

Then, run the following:

```
kubectl apply -f kubernetes/testapp.yaml
```

Once the deployment is running, you can find the publish port by running:

```
kubectl describe service voltestservice
```

Now you can kick the container around and exercise storage using the container's lightweight API.

## Container API

The container API has a few methods, all available via HTTP GET statements. This is NOT a rest API, just something simple for testing.

# /resetfilecheck

Resets the test data to begin a clean test run

# /runfilecheck

Creates the datafiles needed to perform volume function tests

# /textcheck

Returns "1" if the test textfile contains the correct, known data. Returns "0" if not, or if the file doesn't exist.

# /bincheck

Returns "1" if the test binaryfile matches its original checksum. Returns "0" if not (including if the file does not exist)

# /status

Returns "OK" as a container healthcheck

# /shutdown

Immediately terminates the container process. This will trigger Kubernetes to spawn a replacement container.


## TODO items

* Finish test workflow
* Wrap test workflow in actual test assertions
* Gather data needed to configure test environment for Docker storage contributors


## Cheat sheet commands:

Watch container status/restart loop:
kubectl get po -l app=voltest -w

Upgrade container in running deployment:

kubectl patch statefulset voltest --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/image", "value":"khudgins/volcheck:latest"}]'

kubectl describe service voltest
kubectl describe statefulset voltest

./voltestkube --config="/full/path/to/kube.yml"
