package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/drakkan/sftpgo/v2/sdk/plugin/eventsearcher"
	"github.com/hashicorp/go-plugin"
	"github.com/urfave/cli/v2"

	"github.com/sftpgo/sftpgo-plugin-eventsearch/db"
	"github.com/sftpgo/sftpgo-plugin-eventsearch/logger"
)

const (
	version   = "1.0.0-dev"
	envPrefix = "SFTPGO_PLUGIN_EVENTSEARCH_"
)

var (
	commitHash = ""
	buildDate  = ""
)

var (
	driver     string
	instanceID string
	dsn        string

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
				Action: func(c *cli.Context) error {
					logger.AppLogger.Info("starting sftpgo-plugin-eventsearch", "version", getVersionString(),
						"database driver", driver, "instance id", instanceID)
					if err := db.Initialize(driver, dsn); err != nil {
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
