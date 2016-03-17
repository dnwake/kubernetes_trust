package daemon

import (
	// Importing packages here only to make sure their init gets called and
	// therefore they register themselves to the logdriver factory.
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/logger/awslogs"
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/logger/fluentd"
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/logger/gelf"
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/logger/journald"
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/logger/jsonfilelog"
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/logger/splunk"
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/logger/syslog"
)
