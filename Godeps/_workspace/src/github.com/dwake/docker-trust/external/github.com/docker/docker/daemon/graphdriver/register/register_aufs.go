// +build !exclude_graphdriver_aufs,linux

package register

import (
	// register the aufs graphdriver
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/graphdriver/aufs"
)
