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
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	driver := os.Getenv("SFTPGO_PLUGIN_EVENTSEARCH_DRIVER")
	dsn := os.Getenv("SFTPGO_PLUGIN_EVENTSEARCH_DSN")
	if driver == "" || dsn == "" {
		fmt.Println("Driver and/or DSN not set, unable to execute test")
		os.Exit(1)
	}
	if err := Initialize(driver, dsn); err != nil {
		fmt.Printf("unable to initialize database: %v\n", err)
		os.Exit(1)
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	err := sess.AutoMigrate(&providerEventV4{}, &fsEventV5{}, &logEventV1{})
	if err != nil {
		fmt.Printf("unable to migrate database: %v\n", err)
		os.Exit(1)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

type fsEventV5 struct {
	ID                string `gorm:"primaryKey;size:36"`
	Timestamp         int64  `gorm:"size:64;not null;index:idx_fs_events_timestamp"`
	Action            string `gorm:"size:60;not null;index:idx_fs_events_action"`
	Username          string `gorm:"size:255;not null;index:idx_fs_events_username"`
	FsPath            string
	FsTargetPath      string
	VirtualPath       string
	VirtualTargetPath string
	SSHCmd            string `gorm:"size:60;index:idx_fs_events_ssh_cmd"`
	FileSize          int64  `gorm:"size:64"`
	Elapsed           int64  `gorm:"size:64"`
	Status            int    `gorm:"size:32;index:idx_fs_events_status"`
	Protocol          string `gorm:"size:30;not null;index:idx_fs_events_protocol"`
	SessionID         string `gorm:"size:100"`
	IP                string `gorm:"size:50;index:idx_fs_events_ip"`
	FsProvider        int    `gorm:"size:32;index:idx_fs_events_provider"`
	Bucket            string `gorm:"size:512;index:idx_fs_events_bucket"`
	Endpoint          string `gorm:"size:512;index:idx_fs_events_endpoint"`
	OpenFlags         int    `gorm:"size:32"`
	Role              string `gorm:"size:255;index:idx_fs_events_role"`
	InstanceID        string `gorm:"size:60;index:idx_fs_events_instance_id"`
}

func (ev *fsEventV5) TableName() string {
	return "eventstore_fs_events"
}

type providerEventV4 struct {
	ID         string `gorm:"primaryKey;size:36"`
	Timestamp  int64  `gorm:"size:64;not null;index:idx_provider_events__timestamp"`
	Action     string `gorm:"size:60;not null;index:idx_provider_events_action"`
	Username   string `gorm:"size:255;not null;index:idx_provider_events_username"`
	IP         string `gorm:"size:50;index:idx_provider_events_ip"`
	ObjectType string `gorm:"size:50;index:idx_provider_events_object_type"`
	ObjectName string `gorm:"size:255;index:idx_provider_events_object_name"`
	ObjectData []byte
	Role       string `gorm:"size:255;index:idx_provider_events_role"`
	InstanceID string `gorm:"size:60;index:idx_provider_events_instance_id"`
}

func (ev *providerEventV4) TableName() string {
	return "eventstore_provider_events"
}

type logEventV1 struct {
	ID         string `gorm:"primaryKey;size:36"`
	Timestamp  int64  `gorm:"size:64;not null;index:idx_log_events_timestamp"`
	Event      int    `gorm:"size:32;not null;index:idx_log_events_event"`
	Protocol   string `gorm:"size:30;not null;index:idx_log_events_protocol"`
	Username   string `gorm:"size:255;index:idx_log_events_username"`
	IP         string `gorm:"size:50;index:idx_log_events_ip"`
	Message    string `gorm:"size:2048"`
	Role       string `gorm:"size:255;index:idx_log_events_role"`
	InstanceID string `gorm:"size:60;index:idx_log_events_instance_id"`
}

func (ev *logEventV1) TableName() string {
	return "eventstore_log_events"
}
