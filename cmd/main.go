/*
Copyright © 2019 The controller101 Authors.

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
	"context"
	goflag "flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloud-native-taiwan/controller101/pkg/controller"
	cloudnative "github.com/cloud-native-taiwan/controller101/pkg/generated/clientset/versioned"
	cloudnativeinformer "github.com/cloud-native-taiwan/controller101/pkg/generated/informers/externalversions"
	"github.com/cloud-native-taiwan/controller101/pkg/version"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

const defaultSyncTime = time.Second * 30

var (
	kubeconfig  string
	showVersion bool
	threads     int
)

func parseFlags() {
	flag.StringVarP(&kubeconfig, "kubeconfig", "", "", "Absolute path to the kubeconfig file.")
	flag.IntVarP(&threads, "threads", "", 2, "Number of worker threads used by the controller.")
	flag.BoolVarP(&showVersion, "version", "", false, "Display the version.")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}

func restConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	parseFlags()

	if showVersion {
		fmt.Fprintf(os.Stdout, "%s\n", version.GetVersion())
		os.Exit(0)
	}

	k8scfg, err := restConfig(kubeconfig)
	if err != nil {
		klog.Fatalf("Error to build rest config: %s", err.Error())
	}

	clientset, err := cloudnative.NewForConfig(k8scfg)
	if err != nil {
		klog.Fatalf("Error to build cloudnative clientset: %s", err.Error())
	}

	informer := cloudnativeinformer.NewSharedInformerFactory(clientset, defaultSyncTime)
	controller := controller.New(clientset, informer)
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := controller.Run(ctx, threads); err != nil {
		klog.Fatalf("Error to run the controller instance: %s.", err)
	}

	<-signalChan
	cancel()
	controller.Stop()
}
