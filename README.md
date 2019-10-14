# Controller 101
This repository implements a simple controller for watching `VM` resources as defined with a [CRD](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/).

This example will show you how to perform basic operations such as:

* [x] How to create a custom resource of type `VM` using CRD API.
* [x] How to operate instances of type `VM`.
* [x] How to implement a controller for handling an instance of type `VM` to move the current state towards the desired state.
* [x] How to use Finalizer on instances of type `VM`.
* [x] How to implement LeaseLock for multiple controllers.
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
$ eval $(minikube docker-env)
$ POD_NAME=test1 go run cmd/main.go --kubeconfig=$HOME/.kube/config \
    -v=3 --logtostderr \
    --lease-lock-namespace=default \
    --vm-driver=docker
...
I1015 02:16:08.067269   53517 leaderelection.go:242] attempting to acquire leader lease  default/controller101...
I1015 02:16:08.083723   53517 leaderelection.go:252] successfully acquired lease default/controller101
I1015 02:16:08.083830   53517 controller.go:77] Starting the controller
I1015 02:16:08.083846   53517 controller.go:78] Waiting for the informer caches to sync
I1015 02:16:08.185334   53517 controller.go:86] Started workers
I1015 02:16:08.185379   53517 main.go:144] test1: leading
```
