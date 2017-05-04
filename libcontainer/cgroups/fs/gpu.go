// +build linux

package fs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/configs"
)

const (
	cgroupGPUMemoryLimit = "gpu.max_mem_in_bytes"
	cgroupGPUPriority    = "gpu.priority"
)

type GPUGroup struct {
}

func (s *GPUGroup) Name() string {
	return "gpu"
}

func (s *GPUGroup) Apply(d *cgroupData) (err error) {
	path, err := d.path("gpu")
	if err != nil && !cgroups.IsNotFound(err) {
		return err
	} else if path == "" {
		return nil
	}
	fmt.Printf("gpu cgroup %s %d\n", path, d.pid)

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	if err := cgroups.WriteCgroupProc(path, d.pid); err != nil {
		return err
	}

	return nil
}

func (s *GPUGroup) Set(path string, cgroup *configs.Cgroup) error {

	if cgroup.Resources.GPUMemory != 0 {
		if err := writeFile(path, cgroupGPUMemoryLimit, strconv.FormatUint(cgroup.Resources.GPUMemory, 10)); err != nil {
			return err
		}
	}

	if cgroup.Resources.GPUPriority != 0 {
		if err := writeFile(path, cgroupGPUPriority, strconv.FormatInt(cgroup.Resources.GPUPriority, 10)); err != nil {
			return err
		}
	}

	return nil
}

func (s *GPUGroup) Remove(d *cgroupData) error {
	return removePath(d.path("gpu"))
}

func (s *GPUGroup) GetStats(path string, stats *cgroups.Stats) error {
	return nil
}

func gpuAssigned(cgroup *configs.Cgroup) bool {
	return cgroup.Resources.GPUMemory != 0 ||
		cgroup.Resources.GPUPriority != 0
}
