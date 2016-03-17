package daemon

import (
	derr "github.com/dwake/docker-trust/external/github.com/docker/docker/errors"
	"github.com/dwake/docker-trust/external/github.com/docker/engine-api/types"
)

// ContainerTop is not supported on Windows and returns an error.
func (daemon *Daemon) ContainerTop(name string, psArgs string) (*types.ContainerProcessList, error) {
	return nil, derr.ErrorCodeNoTop
}
