# Controller 101
This repository implements a simple controller for watching `VM` resources as defined with a [CRD](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/).

This example will show you how to perform basic operations such as:

* [x] How to create a custom resource of type `VM` using CRD API.
* [ ] How to operate instances of type `VM`.
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
(TBD)