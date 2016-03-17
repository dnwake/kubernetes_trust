// +build linux,!seccomp

package native

import "github.com/dwake/docker-trust/external/github.com/opencontainers/runc/libcontainer/configs"

var (
	defaultSeccompProfile *configs.Seccomp
)
