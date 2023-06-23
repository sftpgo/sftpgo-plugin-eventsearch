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

// LogEvent defines a log event
type LogEvent struct {
	ID         string `json:"id" gorm:"primaryKey"`
	Timestamp  int64  `json:"timestamp"`
	Event      int    `json:"event"`
	Protocol   string `json:"protocol,omitempty"`
	Username   string `json:"username,omitempty"`
	IP         string `json:"ip,omitempty"`
	Message    string `json:"message,omitempty"`
	Role       string `json:"role,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`
}

// TableName defines the database table name
func (ev *LogEvent) TableName() string {
	return "eventstore_log_events"
}
