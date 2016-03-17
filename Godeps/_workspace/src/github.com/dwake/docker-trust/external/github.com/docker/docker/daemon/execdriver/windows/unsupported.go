// +build !windows

package windows

import (
	"fmt"

	"github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/execdriver"
)

// NewDriver returns a new execdriver.Driver
func NewDriver(root, initPath string) (execdriver.Driver, error) {
	return nil, fmt.Errorf("Windows driver not supported on non-Windows")
}
