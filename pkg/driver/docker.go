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

package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	repository = "docker.io/library"
	imageName  = "nginx:1.17.4"
)

type DockerDriver struct {
	client *client.Client
}

func NewDockerDriver() (*DockerDriver, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &DockerDriver{client: cli}, nil
}

func (d *DockerDriver) pullImage(name string) error {
	opts := types.ImageListOptions{Filters: filters.NewArgs()}
	opts.Filters.Add("reference", name)
	imageList, err := d.client.ImageList(context.Background(), opts)
	if err != nil {
		return err
	}

	imageFullName := fmt.Sprintf("%s/%s", repository, name)
	if len(imageList) == 0 {
		_, err := d.client.ImagePull(context.Background(), imageFullName, types.ImagePullOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DockerDriver) CreateServer(req *CreateRequest) (*CreateReponse, error) {
	if err := d.pullImage(imageName); err != nil {
		return nil, err
	}

	config := &container.Config{Image: imageName}
	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			NanoCPUs: req.CPU,
			Memory:   req.Memory,
		},
	}
	resp, err := d.client.ContainerCreate(context.Background(), config, hostConfig, nil, req.Name)
	if err != nil {
		return nil, err
	}

	if err := d.client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	return &CreateReponse{ID: resp.ID}, nil
}

func (d *DockerDriver) GetServerStatus(name string) (*GetStatusReponse, error) {
	status := &GetStatusReponse{}
	container, err := d.getContainer(name)
	if err != nil {
		return nil, err
	}
	status.State = container.State

	resp, err := d.client.ContainerStats(context.Background(), name, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	v := &types.StatsJSON{}
	if err := json.Unmarshal(content, &v); err != nil {
		return nil, err
	}

	var memPercent = 0.0
	previousCPU := v.PreCPUStats.CPUUsage.TotalUsage
	previousSystem := v.PreCPUStats.SystemUsage
	cpuPercent := calculateCPUPercentUnix(previousCPU, previousSystem, v)

	if v.MemoryStats.Limit != 0 {
		memPercent = float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100.0
	}

	status.CPUPercentage = cpuPercent
	status.MemoryPercentage = memPercent
	return status, nil
}

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

func (d *DockerDriver) DeleteServer(name string) error {
	if err := d.client.ContainerStop(context.Background(), name, nil); err != nil {
		return err
	}
	return d.client.ContainerRemove(context.Background(), name, types.ContainerRemoveOptions{})
}

func (d *DockerDriver) IsServerExist(name string) (bool, error) {
	container, err := d.getContainer(name)
	if container != nil {
		return true, nil
	}
	return false, err
}

func (d *DockerDriver) getContainer(name string) (*types.Container, error) {
	opts := types.ContainerListOptions{Filters: filters.NewArgs()}
	opts.Filters.Add("name", name)
	containers, err := d.client.ContainerList(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	if len(containers) < 1 || len(containers) > 1 {
		return nil, fmt.Errorf("Failed to get container by '%s' name", name)
	}
	return &containers[0], nil
}
