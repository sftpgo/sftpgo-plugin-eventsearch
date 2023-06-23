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
	"encoding/json"
	"testing"

	"github.com/rs/xid"
	"github.com/sftpgo/sdk/plugin/eventsearcher"
	"github.com/stretchr/testify/assert"
)

func TestSearchFsEvents(t *testing.T) {
	fsEvents := []FsEvent{
		{
			ID:                xid.New().String(),
			Timestamp:         100,
			Action:            "upload",
			Username:          "username1",
			FsPath:            "/tmp/file.txt",
			FsTargetPath:      "/tmp/target.txt",
			VirtualPath:       "file.txt",
			VirtualTargetPath: "target.txt",
			SSHCmd:            "scp",
			FileSize:          123,
			Status:            1,
			Protocol:          "SFTP",
			IP:                "::1",
			SessionID:         "1",
			InstanceID:        "instance1",
			FsProvider:        0,
			Bucket:            "",
			Endpoint:          "endpoint1",
			OpenFlags:         512,
			Role:              "role1",
		},
		{
			ID:                xid.New().String(),
			Timestamp:         101,
			Action:            "download",
			Username:          "username2",
			FsPath:            "/tmp/file.txt",
			FsTargetPath:      "/tmp/target.txt",
			VirtualPath:       "file.txt",
			VirtualTargetPath: "target.txt",
			SSHCmd:            "",
			FileSize:          123,
			Status:            2,
			Protocol:          "SFTP",
			IP:                "::1",
			SessionID:         "2",
			InstanceID:        "instance2",
			FsProvider:        0,
			Bucket:            "",
			Endpoint:          "",
			OpenFlags:         0,
			Role:              "role1",
		},
		{
			ID:                xid.New().String(),
			Timestamp:         101,
			Action:            "upload",
			Username:          "username3",
			FsPath:            "/tmp/file.txt",
			FsTargetPath:      "/tmp/target.txt",
			VirtualPath:       "file.txt",
			VirtualTargetPath: "target.txt",
			SSHCmd:            "",
			FileSize:          123,
			Status:            2,
			Protocol:          "SFTP",
			IP:                "::1",
			SessionID:         "3",
			InstanceID:        "instance1",
			FsProvider:        0,
			Bucket:            "",
			Endpoint:          "",
			OpenFlags:         0,
			Role:              "role2",
		},
		{
			ID:                xid.New().String(),
			Timestamp:         101,
			Action:            "download",
			Username:          "username4",
			FsPath:            "/tmp/file.txt",
			FsTargetPath:      "/tmp/target.txt",
			VirtualPath:       "file.txt",
			VirtualTargetPath: "target.txt",
			SSHCmd:            "scp",
			FileSize:          123,
			Status:            2,
			Protocol:          "SFTP",
			IP:                "::1",
			SessionID:         "4",
			InstanceID:        "instance2",
			FsProvider:        2,
			Bucket:            "bucket4",
			Endpoint:          "",
			OpenFlags:         0,
		},
		{
			ID:                xid.New().String(),
			Timestamp:         102,
			Action:            "download",
			Username:          "username5",
			FsPath:            "/tmp/file.txt",
			FsTargetPath:      "/tmp/target.txt",
			VirtualPath:       "file.txt",
			VirtualTargetPath: "target.txt",
			SSHCmd:            "scp",
			FileSize:          123,
			Status:            3,
			Protocol:          "SCP",
			IP:                "127.0.0.1",
			SessionID:         "5",
			InstanceID:        "instance3",
			FsProvider:        1,
			Bucket:            "bucket5",
			Endpoint:          "",
			OpenFlags:         0,
		},
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	err := sess.Create(&fsEvents).Error
	assert.NoError(t, err)

	s := Searcher{}
	_, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 0,
		},
	})
	assert.ErrorIs(t, err, errNoLimit)

	// test order ASC
	data, err := s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 1,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	var events []FsEvent
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	assert.Equal(t, fsEvents[0].ID, events[0].ID)
	assert.Equal(t, fsEvents[4].ID, events[4].ID)
	// test order DESC
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 0,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	assert.Equal(t, fsEvents[4].ID, events[0].ID)
	assert.Equal(t, fsEvents[0].ID, events[4].ID)
	// test limit and pagination
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 2,
			Order: 1,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[0].ID, events[0].ID)
	assert.Equal(t, fsEvents[1].ID, events[1].ID)
	// get next page
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			StartTimestamp: events[1].Timestamp,
			Limit:          2,
			Order:          1,
			FromID:         events[1].ID,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[2].ID, events[0].ID)
	assert.Equal(t, fsEvents[3].ID, events[1].ID)
	// get last page
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			StartTimestamp: events[1].Timestamp,
			Limit:          2,
			Order:          1,
			FromID:         events[1].ID,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, fsEvents[4].ID, events[0].ID)
	// get previous page
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			EndTimestamp: events[0].Timestamp,
			Limit:        2,
			Order:        0,
			FromID:       events[0].ID,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[3].ID, events[0].ID)
	assert.Equal(t, fsEvents[2].ID, events[1].ID)
	// get first page
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			EndTimestamp: events[1].Timestamp,
			Limit:        2,
			Order:        0,
			FromID:       events[1].ID,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[1].ID, events[0].ID)
	assert.Equal(t, fsEvents[0].ID, events[1].ID)
	// paginate starting from DESC
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 2,
			Order: 0,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[4].ID, events[0].ID)
	assert.Equal(t, fsEvents[3].ID, events[1].ID)
	// get next page
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			EndTimestamp: events[1].Timestamp,
			Limit:        2,
			Order:        0,
			FromID:       events[1].ID,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[2].ID, events[0].ID)
	assert.Equal(t, fsEvents[1].ID, events[1].ID)
	// get last page
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			EndTimestamp: events[1].Timestamp,
			Limit:        2,
			Order:        0,
			FromID:       events[1].ID,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, fsEvents[0].ID, events[0].ID)
	// get previous page
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			StartTimestamp: events[0].Timestamp,
			Limit:          2,
			Order:          1,
			FromID:         events[0].ID,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[1].ID, events[0].ID)
	assert.Equal(t, fsEvents[2].ID, events[1].ID)
	// get first page
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			StartTimestamp: events[1].Timestamp,
			Limit:          2,
			Order:          1,
			FromID:         events[1].ID,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[3].ID, events[0].ID)
	assert.Equal(t, fsEvents[4].ID, events[1].ID)

	// test other search conditions
	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Username: "username1",
			Limit:    100,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, fsEvents[0], events[0])

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 1,
			IP:    "::1",
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 4)
	assert.Equal(t, fsEvents[:4], events)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		SSHCmd:     "scp",
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		SSHCmd:     "sha256sum",
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		Actions:    []string{"upload", "download", "rename"},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		Actions:    []string{"rename"},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		Protocols:  []string{"SFTP", "HTTP"},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 4)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		Protocols:  []string{"SCP"},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			InstanceIDs: []string{"instance1"},
			Limit:       100,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			InstanceIDs: []string{"instance1", "instance3"},
			Limit:       100,
		},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			InstanceIDs: []string{"instance1", "instance2", "instance3"},
			Limit:       100,
		},
		Statuses:   []int32{1},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			InstanceIDs: []string{"instance1", "instance2", "instance3"},
			Limit:       100,
		},
		Statuses:   []int32{1, 2, 3},
		FsProvider: -1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		FsProvider: 0,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		FsProvider: 1,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		FsProvider: 100,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		FsProvider: -1,
		Bucket:     "bucket5",
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		FsProvider: -1,
		Bucket:     "bucke",
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		FsProvider: -1,
		Endpoint:   "endpoint1",
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		FsProvider: -1,
		Endpoint:   "endpo",
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Role:  "role1",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Role:  "role2",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchFsEvents(&eventsearcher.FsEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Role:  "role3",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	err = sess.Delete(&fsEvents).Error
	assert.NoError(t, err)
}

func TestSearchProviderEvents(t *testing.T) {
	providerEvents := []ProviderEvent{
		{
			ID:         xid.New().String(),
			Timestamp:  100,
			Action:     "add",
			Username:   "username1",
			IP:         "127.1.1.1",
			ObjectType: "api_key",
			ObjectName: "123",
			ObjectData: []byte("data"),
			Role:       "role1",
			InstanceID: "instance1",
		},
		{
			ID:         xid.New().String(),
			Timestamp:  101,
			Action:     "delete",
			Username:   "username2",
			IP:         "127.1.0.1",
			ObjectType: "admin",
			ObjectName: "456",
			ObjectData: []byte("data"),
			InstanceID: "instance2",
		},
		{
			ID:         xid.New().String(),
			Timestamp:  101,
			Action:     "update",
			Username:   "username3",
			IP:         "127.1.0.1",
			ObjectType: "user",
			ObjectName: "678",
			ObjectData: []byte("data"),
			InstanceID: "instance1",
		},
		{
			ID:         xid.New().String(),
			Timestamp:  101,
			Action:     "update",
			Username:   "username4",
			IP:         "127.1.0.1",
			ObjectType: "user",
			ObjectName: "678",
			ObjectData: []byte("data"),
			InstanceID: "instance1",
		},
		{
			ID:         xid.New().String(),
			Timestamp:  102,
			Action:     "update",
			Username:   "username5",
			IP:         "127.1.0.1",
			ObjectType: "admin",
			ObjectName: "0123",
			ObjectData: []byte("data"),
			InstanceID: "instance3",
		},
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	err := sess.Create(&providerEvents).Error
	assert.NoError(t, err)

	s := Searcher{}
	_, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{})
	assert.ErrorIs(t, err, errNoLimit)
	// test order ASC
	data, err := s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 1,
		},
	})
	assert.NoError(t, err)
	var events []ProviderEvent
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	assert.Equal(t, providerEvents[0].ID, events[0].ID)
	assert.Equal(t, providerEvents[4].ID, events[4].ID)
	for _, ev := range events {
		assert.Equal(t, []byte("data"), ev.ObjectData)
	}
	// test omit object data
	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 1,
			Order: 1,
		},
		OmitObjectData: true,
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	if assert.Len(t, events, 1) {
		assert.Nil(t, events[0].ObjectData)
	}
	// test order DESC
	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 0,
		},
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	assert.Equal(t, providerEvents[4].ID, events[0].ID)
	assert.Equal(t, providerEvents[0].ID, events[4].ID)
	// test limit and pagination
	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 3,
			Order: 1,
		},
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
	assert.Equal(t, providerEvents[0].ID, events[0].ID)
	assert.Equal(t, providerEvents[2].ID, events[2].ID)
	// get next page
	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:          3,
			Order:          1,
			StartTimestamp: events[2].Timestamp,
			FromID:         events[2].ID,
		},
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, providerEvents[3].ID, events[0].ID)
	assert.Equal(t, providerEvents[4].ID, events[1].ID)
	// get previous page
	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:        3,
			Order:        0,
			EndTimestamp: events[0].Timestamp,
			FromID:       events[0].ID,
		},
	})
	assert.NoError(t, err)

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
	assert.Equal(t, providerEvents[2].ID, events[0].ID)
	assert.Equal(t, providerEvents[0].ID, events[2].ID)
	// test other search conditions
	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 0,
		},
		Actions: []string{"add", "delete"},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:    100,
			Order:    0,
			Username: "username1",
			IP:       "127.1.1.1",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, providerEvents[0].ID, events[0].ID)

	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 0,
		},
		ObjectTypes: []string{"api_key", "user"},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)

	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 0,
		},
		ObjectName:  "123",
		ObjectTypes: []string{"api_key", "user"},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, providerEvents[0].ID, events[0].ID)

	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:       100,
			Order:       0,
			InstanceIDs: []string{"instance2", "instance3"},
		},
		ObjectName:  "123",
		ObjectTypes: []string{"api_key", "user"},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:       100,
			Order:       0,
			InstanceIDs: []string{"instance2", "instance3"},
		},
		ObjectTypes: []string{"api_key", "admin"},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Role:  "role1",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchProviderEvents(&eventsearcher.ProviderEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Role:  "role3",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	err = sess.Delete(&providerEvents).Error
	assert.NoError(t, err)
}

func TestSearchLogEvents(t *testing.T) {
	logEvents := []LogEvent{
		{
			ID:         xid.New().String(),
			Timestamp:  100,
			Event:      1,
			Protocol:   "SSH",
			Username:   "username1",
			IP:         "127.1.1.1",
			Message:    "error1",
			Role:       "role1",
			InstanceID: "instance1",
		},
		{
			ID:         xid.New().String(),
			Timestamp:  101,
			Event:      2,
			Protocol:   "FTP",
			Username:   "username2",
			IP:         "127.1.0.1",
			Message:    "error2",
			InstanceID: "instance2",
		},
		{
			ID:         xid.New().String(),
			Timestamp:  101,
			Event:      3,
			Protocol:   "FTP",
			Username:   "username3",
			IP:         "127.1.0.1",
			Message:    "error3",
			InstanceID: "instance1",
		},
		{
			ID:         xid.New().String(),
			Timestamp:  101,
			Event:      3,
			Protocol:   "HTTP",
			Username:   "username4",
			IP:         "127.1.0.1",
			Message:    "error4",
			InstanceID: "instance1",
		},
		{
			ID:         xid.New().String(),
			Timestamp:  102,
			Event:      2,
			Protocol:   "DAV",
			Username:   "username5",
			IP:         "127.1.0.1",
			Message:    "error5",
			InstanceID: "instance3",
		},
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	err := sess.Create(&logEvents).Error
	assert.NoError(t, err)

	s := Searcher{}
	_, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{})
	assert.ErrorIs(t, err, errNoLimit)
	// test order ASC
	data, err := s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 1,
		},
	})
	assert.NoError(t, err)
	var events []LogEvent
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	assert.Equal(t, logEvents[0].ID, events[0].ID)
	assert.Equal(t, logEvents[4].ID, events[4].ID)
	// test order DESC
	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Order: 0,
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	assert.Equal(t, logEvents[4].ID, events[0].ID)
	assert.Equal(t, logEvents[0].ID, events[4].ID)
	// test limit and pagination
	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 3,
			Order: 1,
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
	assert.Equal(t, logEvents[0].ID, events[0].ID)
	assert.Equal(t, logEvents[2].ID, events[2].ID)
	// get next page
	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:          3,
			Order:          1,
			StartTimestamp: events[2].Timestamp,
			FromID:         events[2].ID,
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, logEvents[3].ID, events[0].ID)
	assert.Equal(t, logEvents[4].ID, events[1].ID)
	// get previous page
	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:        3,
			Order:        0,
			EndTimestamp: events[0].Timestamp,
			FromID:       events[0].ID,
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
	assert.Equal(t, logEvents[2].ID, events[0].ID)
	assert.Equal(t, logEvents[0].ID, events[2].ID)
	// test other search conditions
	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		Events: []int32{100, 101},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		Events: []int32{1},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
		},
		Protocols: []string{"FTP", "DAV"},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:    100,
			Username: "u1",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:    100,
			Username: "username1",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:    100,
			Username: "username1",
			IP:       "127.1.1.1",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:    100,
			Username: "username1",
			IP:       "127.1.0.1",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit:       100,
			InstanceIDs: []string{"instance2", "instance3"},
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Role:  "role1",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, err = s.SearchLogEvents(&eventsearcher.LogEventSearch{
		CommonSearchParams: eventsearcher.CommonSearchParams{
			Limit: 100,
			Role:  "role123",
		},
	})
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	err = sess.Delete(&logEvents).Error
	assert.NoError(t, err)
}
