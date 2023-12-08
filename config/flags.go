// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	ConfigFileFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "JSON configuration file",
	}

	VerbosityFlag = &cli.StringFlag{
		Name:  "verbosity",
		Usage: "Supports levels crit (silent) to trce (trace)",
		Value: log.InfoLevel.String(),
	}

	ResetFlag = &cli.BoolFlag{
		Name:  "reset",
		Usage: "Resets the local DB and Rabbit MQ queue and starts from fresh.",
		Value: false,
	}

	LatestBlockFlag = &cli.BoolFlag{
		Name:  "latest",
		Usage: "Overrides blockstore and start block, starts from latest block",
	}

	MetricsFlag = &cli.BoolFlag{
		Name:  "metrics",
		Usage: "Enables metric server",
		Value: true,
	}

	MetricsPort = &cli.StringFlag{
		Name:  "metricsPort",
		Usage: "Port to serve metrics on",
		Value: "8001",
	}
)
