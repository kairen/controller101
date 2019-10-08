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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudnative "github.com/cloud-native-taiwan/controller101/pkg/generated/clientset/versioned"
	cloudnativeinformer "github.com/cloud-native-taiwan/controller101/pkg/generated/informers/externalversions"
	listerv1alpha1 "github.com/cloud-native-taiwan/controller101/pkg/generated/listers/cloudnative/v1alpha1"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

const (
	resouceName = "VirtualMachine"
)

type Controller struct {
	clientset cloudnative.Interface
	informer  cloudnativeinformer.SharedInformerFactory
	lister    listerv1alpha1.VirtualMachineLister
	synced    cache.InformerSynced
	queue     workqueue.RateLimitingInterface
}

func New(clientset cloudnative.Interface, informer cloudnativeinformer.SharedInformerFactory) *Controller {
	vmInformer := informer.Cloudnative().V1alpha1().VirtualMachines()
	controller := &Controller{
		clientset: clientset,
		informer:  informer,
		lister:    vmInformer.Lister(),
		synced:    vmInformer.Informer().HasSynced,
		queue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), resouceName),
	}

	vmInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueue,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueue(new)
		},
	})
	return controller
}

func (c *Controller) Run(ctx context.Context, threadiness int) error {
	go c.informer.Start(ctx.Done())
	klog.Info("Starting the controller")
	klog.Info("Waiting for the informer caches to sync")
	if ok := cache.WaitForCacheSync(ctx.Done(), c.synced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, ctx.Done())
	}
	klog.Info("Started workers")
	return nil
}

func (c *Controller) Stop() {
	glog.Info("Stopping the controller")
	c.queue.ShutDown()
}

func (c *Controller) runWorker() {
	defer utilruntime.HandleCrash()
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.queue.Done(obj)
		key, ok := obj.(string)
		if !ok {
			c.queue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("Controller expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.queue.AddRateLimited(key)
			return fmt.Errorf("Controller error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.queue.Forget(obj)
		glog.Infof("Controller successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}
	return true
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.queue.Add(key)
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return err
	}

	vm, err := c.lister.VirtualMachines(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("virtualmachine '%s' in work queue no longer exists", key))
			return err
		}
		return err
	}

	data, err := json.Marshal(vm)
	if err != nil {
		return err
	}

	klog.Infof("Controller get %s/%s object: %s", namespace, name, string(data))
	return nil
}
