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

// ProviderEvent defines a provider event
type ProviderEvent struct {
	ID         string `json:"id" gorm:"primaryKey"`
	Timestamp  int64  `json:"timestamp"`
	Action     string `json:"action"`
	Username   string `json:"username"`
	IP         string `json:"ip,omitempty"`
	ObjectType string `json:"object_type"`
	ObjectName string `json:"object_name"`
	ObjectData []byte `json:"object_data,omitempty"`
	Role       string `json:"role,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`
}

// TableName defines the database table name
func (ev *ProviderEvent) TableName() string {
	return "eventstore_provider_events"
}
