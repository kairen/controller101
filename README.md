# Controller 101
This repository implements a simple controller for watching `VM` resources as defined with a [CRD](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/).

This example will show you how to perform basic operations such as:

* [x] How to create a custom resource of type `VM` using CRD API.
* [x] How to operate instances of type `VM`.
* [ ] How to implement a controller for handling an instance of type `VM` to move the current state towards the desired state.
* [ ] How to use Finalizer on instances of type `VM`.
* [ ] How to implement LeaseLock for multiple controllers.
* [ ] How to expose metrics of the controller.
 
## Building from Source
Clone the repo in whatever working directory you like, and run the following commands:

```sh
$ export GO111MODULE=on
$ git clone https://github.com/cloud-native-taiwan/controller101
$ cd controller101
$ make
```

## Running
Run the following command to debug:

```sh
$ minikube start --kubernetes-version=v1.15.4 
$ go run cmd/main.go --kubeconfig=$HOME/.kube/config -v=2 --logtostderr
I1008 15:38:30.350446   52017 controller.go:68] Starting the controller
I1008 15:38:30.350543   52017 controller.go:69] Waiting for the informer caches to sync
I1008 15:38:30.454799   52017 controller.go:77] Started workers
```