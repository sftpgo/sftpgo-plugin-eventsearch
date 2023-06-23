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

// FsEvent defines a filesystem event
type FsEvent struct {
	ID                string `json:"id" gorm:"primaryKey"`
	Timestamp         int64  `json:"timestamp"`
	Action            string `json:"action"`
	Username          string `json:"username"`
	FsPath            string `json:"fs_path"`
	FsTargetPath      string `json:"fs_target_path,omitempty"`
	VirtualPath       string `json:"virtual_path"`
	VirtualTargetPath string `json:"virtual_target_path,omitempty"`
	SSHCmd            string `json:"ssh_cmd,omitempty"`
	FileSize          int64  `json:"file_size,omitempty"`
	Elapsed           int64  `json:"elapsed,omitempty"`
	Status            int    `json:"status"`
	Protocol          string `json:"protocol"`
	IP                string `json:"ip,omitempty"`
	SessionID         string `json:"session_id"`
	FsProvider        int    `json:"fs_provider"`
	Bucket            string `json:"bucket,omitempty"`
	Endpoint          string `json:"endpoint,omitempty"`
	OpenFlags         int    `json:"open_flags,omitempty"`
	Role              string `json:"role,omitempty"`
	InstanceID        string `json:"instance_id,omitempty"`
}

// TableName defines the database table name
func (ev *FsEvent) TableName() string {
	return "eventstore_fs_events"
}
