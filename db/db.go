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

package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/sftpgo/sftpgo-plugin-eventsearch/logger"
)

const (
	driverNamePostgreSQL = "postgres"
	driverNameMySQL      = "mysql"
)

var (
	handle              *gorm.DB
	defaultQueryTimeout = 20 * time.Second
)

// Initialize initializes the database engine
func Initialize(driver, dsn, customTLSConfig string) error {
	var err error

	switch driver {
	case driverNamePostgreSQL:
		handle, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 gormlogger.Discard,
		})
		if err != nil {
			logger.AppLogger.Error("unable to create db handle", "error", err)
			return err
		}
	case driverNameMySQL:
		if err := handleCustomTLSConfig(customTLSConfig); err != nil {
			logger.AppLogger.Error("unable to register custom tls config", "error", err)
			return err
		}
		handle, err = gorm.Open(mysql.New(mysql.Config{
			DSN: dsn,
		}), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 gormlogger.Discard,
		})
		if err != nil {
			logger.AppLogger.Error("unable to create db handle", "error", err)
			return err
		}
	default:
		return fmt.Errorf("unsupported database driver %v", driver)
	}

	sqlDB, err := handle.DB()
	if err != nil {
		logger.AppLogger.Error("unable to get sql db handle", "error", err)
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(4 * time.Minute)
	sqlDB.SetConnMaxLifetime(2 * time.Minute)

	return sqlDB.Ping()
}

// getDefaultSession returns a database session with the default timeout.
// Don't forget to cancel the returned context
func getDefaultSession() (*gorm.DB, context.CancelFunc) {
	return getSessionWithTimeout(defaultQueryTimeout)
}

// getSessionWithTimeout returns a database session with the specified timeout.
// Don't forget to cancel the returned context
func getSessionWithTimeout(timeout time.Duration) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	return handle.WithContext(ctx), cancel
}

func handleCustomTLSConfig(config string) error {
	if config == "" {
		return nil
	}
	values, err := url.ParseQuery(config)
	if err != nil {
		logger.AppLogger.Error("unable to parse custom tls config", "value", config, "error", err)
		return fmt.Errorf("unable to parse tls config: %w", err)
	}
	rootCert := values.Get("root_cert")
	clientCert := values.Get("client_cert")
	clientKey := values.Get("client_key")
	tlsMode := values.Get("tls_mode")

	tlsConfig := &tls.Config{}
	if rootCert != "" {
		rootCAs, err := x509.SystemCertPool()
		if err != nil {
			rootCAs = x509.NewCertPool()
		}
		rootCrt, err := os.ReadFile(rootCert)
		if err != nil {
			return fmt.Errorf("unable to load root certificate %q: %v", rootCert, err)
		}
		if !rootCAs.AppendCertsFromPEM(rootCrt) {
			return fmt.Errorf("unable to parse root certificate %q", rootCert)
		}
		tlsConfig.RootCAs = rootCAs
	}
	if clientCert != "" && clientKey != "" {
		cert := make([]tls.Certificate, 0, 1)
		tlsCert, err := tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			return fmt.Errorf("unable to load key pair %q, %q: %v", clientCert, clientKey, err)
		}
		cert = append(cert, tlsCert)
		tlsConfig.Certificates = cert
	}
	if tlsMode == "1" {
		tlsConfig.InsecureSkipVerify = true
	}

	if err := mysqldriver.RegisterTLSConfig("custom", tlsConfig); err != nil {
		return fmt.Errorf("unable to register tls config: %v", err)
	}
	return nil
}
