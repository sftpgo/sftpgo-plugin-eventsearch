# SFTPGo events search plugin

![Build](https://github.com/sftpgo/sftpgo-plugin-eventsearch/workflows/Build/badge.svg)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPLv3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

This plugin allows to search [SFTPGo](https://github.com/drakkan/sftpgo/) filesystem and provider events stored using the [sftpgo-plugin-eventstore](https://github.com/sftpgo/sftpgo-plugin-eventstore).

## Configuration

The plugin can be configured within the `plugins` section of the SFTPGo configuration file. To start the plugin you have to use the `serve` subcommand. Here is the usage.

```shell
NAME:
   sftpgo-plugin-eventsearch serve - Launch the SFTPGo plugin, it must be called from an SFTPGo instance

USAGE:
   sftpgo-plugin-eventsearch serve [command options]

OPTIONS:
   --driver value      Database driver (required) [$SFTPGO_PLUGIN_EVENTSEARCH_DRIVER]
   --dsn value         Data source URI (required) [$SFTPGO_PLUGIN_EVENTSEARCH_DSN]
   --custom-tls value  Custom TLS config for MySQL driver (optional) [$SFTPGO_PLUGIN_EVENTSEARCH_CUSTOM_TLS]
   --pool-size value   Naximum number of open database connections (default: 0) [$SFTPGO_PLUGIN_EVENTSEARCH_POOL_SIZE]
   --help, -h          show help
```

The `driver` and `dsn` flags are required and must match the ones configured for [sftpgo-plugin-eventstore](https://github.com/sftpgo/sftpgo-plugin-eventstore).
Each flag can also be set using environment variables, for example the DSN can be set using the `SFTPGO_PLUGIN_EVENTSEARCH_DSN` environment variable.

This is an example configuration.

```json
...
"plugins": [
    {
      "type": "eventsearcher",
      "cmd": "<path to sftpgo-plugin-eventsearch>",
      "args": ["serve", "--driver", "postgres"],
      "sha256sum": "",
      "auto_mtls": true
    }
  ]
...
```

With the above example the plugin is configured to connect to PostgreSQL. We set the DSN using the `SFTPGO_PLUGIN_EVENTSEARCH_DSN` environment variable. You can now search events in SFTPGo using REST API or the WebAdmin UI.

## Supported database services

### PostgreSQL

To use Postgres you have to use `postgres` as driver. If you have a database named `sftpgo_events` on localhost and you want to connect to it using the user `sftpgo` with the password `sftpgopass` you can use a DSN like the following one.

```shell
"host='127.0.0.1' port=5432 dbname='sftpgo_events' user='sftpgo' password='sftpgopass' sslmode=disable connect_timeout=10"
```

Please refer to the documentation [here](https://github.com/go-gorm/postgres) for details about the dsn.

### MariaDB/MySQL

To use MariaDB/MySQL you have to use `mysql` as driver. If you have a database named `sftpgo_events` on localhost and you want to connect to it using the user `sftpgo` with the password `sftpgopass` you can use a DSN like the following one.

```shell
"sftpgo:sftpgopass@tcp([127.0.0.1]:3306)/sftpgo_events?collation=utf8mb4_unicode_ci&interpolateParams=true&timeout=10s&tls=false&writeTimeout=30s&readTimeout=30s&parseTime=true&clientFoundRows=true"
```

Please refer to the documentation [here](https://github.com/go-gorm/mysql) for details about the dsn.
