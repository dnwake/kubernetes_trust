// +build !experimental

package daemon

import flag "github.com/dwake/docker-trust/external/github.com/docker/docker/pkg/mflag"

func (config *Config) attachExperimentalFlags(cmd *flag.FlagSet, usageFn func(string) string) {
}
