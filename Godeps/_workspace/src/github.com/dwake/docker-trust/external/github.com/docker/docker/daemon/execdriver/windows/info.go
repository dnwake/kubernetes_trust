// +build windows

package windows

import (
	"github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/execdriver"
	"github.com/dwake/docker-trust/external/github.com/docker/engine-api/types/container"
)

type info struct {
	ID        string
	driver    *Driver
	isolation container.IsolationLevel
}

// Info implements the exec driver Driver interface.
func (d *Driver) Info(id string) execdriver.Info {
	return &info{
		ID:        id,
		driver:    d,
		isolation: DefaultIsolation,
	}
}

func (i *info) IsRunning() bool {
	var running bool
	running = true // TODO Need an HCS API
	return running
}
