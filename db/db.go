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
	"fmt"
	"time"

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
func Initialize(driver, dsn string) error {
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
