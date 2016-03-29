package cli

import flag "github.com/dwake/docker-trust/external/github.com/docker/docker/pkg/mflag"

// ClientFlags represents flags for the docker client.
type ClientFlags struct {
	FlagSet   *flag.FlagSet
	Common    *CommonFlags
	PostParse func()

	ConfigDir string
}
