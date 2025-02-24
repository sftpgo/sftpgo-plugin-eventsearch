// Copyright (C) 2021-2023 Nicola Murino
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/sftpgo/sdk/plugin/eventsearcher"
	"github.com/urfave/cli/v2"

	"github.com/sftpgo/sftpgo-plugin-eventsearch/db"
	"github.com/sftpgo/sftpgo-plugin-eventsearch/logger"
)

const (
	version   = "1.0.20"
	envPrefix = "SFTPGO_PLUGIN_EVENTSEARCH_"
)

var (
	commitHash = ""
	buildDate  = ""
)

var (
	driver          string
	instanceID      string
	dsn             string
	customTLSConfig string
	poolSize        int

	serveFlags = []cli.Flag{
		&cli.StringFlag{
			Name:        "driver",
			Usage:       "Database driver (required)",
			Destination: &driver,
			EnvVars:     []string{envPrefix + "DRIVER"},
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "dsn",
			Usage:       "Data source URI (required)",
			Destination: &dsn,
			EnvVars:     []string{envPrefix + "DSN"},
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "custom-tls",
			Usage:       "Custom TLS config for MySQL driver (optional)",
			Destination: &customTLSConfig,
			EnvVars:     []string{envPrefix + "CUSTOM_TLS"},
			Required:    false,
		},
		&cli.IntFlag{
			Name:        "pool-size",
			Usage:       "Naximum number of open database connections",
			Destination: &poolSize,
			EnvVars:     []string{envPrefix + "POOL_SIZE"},
			Required:    false,
		},
	}

	rootCmd = &cli.App{
		Name:    "sftpgo-plugin-eventsearch",
		Version: getVersionString(),
		Usage:   "SFTPGo events store plugin",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Launch the SFTPGo plugin, it must be called from an SFTPGo instance",
				Flags: serveFlags,
				Action: func(_ *cli.Context) error {
					logger.AppLogger.Info("starting sftpgo-plugin-eventsearch", "version", getVersionString(),
						"database driver", driver, "instance id", instanceID, "pool size", poolSize)
					if err := db.Initialize(driver, dsn, customTLSConfig, poolSize); err != nil {
						logger.AppLogger.Error("unable to initialize database", "error", err)
						return err
					}

					plugin.Serve(&plugin.ServeConfig{
						HandshakeConfig: eventsearcher.Handshake,
						Plugins: map[string]plugin.Plugin{
							eventsearcher.PluginName: &eventsearcher.Plugin{Impl: &db.Searcher{}},
						},
						GRPCServer: plugin.DefaultGRPCServer,
					})

					return errors.New("the plugin exited unexpectedly")
				},
			},
		},
	}
)

// Execute runs the root command
func Execute() error {
	return rootCmd.Run(os.Args)
}

func getVersionString() string {
	var sb strings.Builder
	sb.WriteString(version)
	if commitHash != "" {
		sb.WriteString("-")
		sb.WriteString(commitHash)
	}
	if buildDate != "" {
		sb.WriteString("-")
		sb.WriteString(buildDate)
	}
	return sb.String()
}
