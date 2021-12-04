package db

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchFsEvents(t *testing.T) {
	fsEvents := []FsEvent{
		{
			ID:                "1",
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
		},
		{
			ID:                "2",
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
		},
		{
			ID:                "3",
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
		},
		{
			ID:                "4",
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
		},
		{
			ID:                "5",
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
		},
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	err := sess.Create(&fsEvents).Error
	assert.NoError(t, err)

	s := Searcher{}
	_, _, _, err = s.SearchFsEvents(0, 0, "", "", "", nil, nil, nil, nil, nil, 0, 0)
	assert.ErrorIs(t, err, errNoLimit)

	// test order ASC
	data, sameAtStart, sameAtEnd, err := s.SearchFsEvents(0, 0, "", "", "", nil, nil, nil, nil, nil, 100, 1)
	assert.NoError(t, err)

	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, fsEvents[0].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 1)
	assert.Equal(t, fsEvents[4].ID, sameAtEnd[0])

	var events []FsEvent
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	// test order DESC
	data, sameAtStart, sameAtEnd, err = s.SearchFsEvents(0, 0, "", "", "", nil, nil, nil, nil, nil, 100, 0)
	assert.NoError(t, err)

	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, fsEvents[4].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 1)
	assert.Equal(t, fsEvents[0].ID, sameAtEnd[0])

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	// test limit and pagination
	data, sameAtStart, sameAtEnd, err = s.SearchFsEvents(0, 0, "", "", "", nil, nil, nil, nil, nil, 3, 1)
	assert.NoError(t, err)
	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, fsEvents[0].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 2)
	assert.Equal(t, fsEvents[1].ID, sameAtEnd[1])
	assert.Equal(t, fsEvents[2].ID, sameAtEnd[0])

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
	// get next page
	data, sameAtStart, sameAtEnd, err = s.SearchFsEvents(events[2].Timestamp, 0, "", "", "", nil, nil, nil, sameAtEnd, nil,
		3, 1)
	assert.NoError(t, err)
	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, fsEvents[3].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 1)
	assert.Equal(t, fsEvents[4].ID, sameAtEnd[0])

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, fsEvents[3].ID, events[0].ID)
	assert.Equal(t, fsEvents[4].ID, events[1].ID)
	// get previous page
	data, sameAtStart, sameAtEnd, err = s.SearchFsEvents(0, events[0].Timestamp, "", "", "", nil, nil, nil, sameAtStart, nil,
		3, 1)
	assert.NoError(t, err)
	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, fsEvents[0].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 2)
	assert.Equal(t, fsEvents[1].ID, sameAtEnd[1])
	assert.Equal(t, fsEvents[2].ID, sameAtEnd[0])

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
	assert.Equal(t, fsEvents[:3], events)
	// test other search conditions
	data, _, _, err = s.SearchFsEvents(0, 0, "username1", "", "", nil, nil, nil, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, fsEvents[0], events[0])

	data, _, _, err = s.SearchFsEvents(0, 0, "", "::1", "", nil, nil, nil, nil, nil, 100, 1)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 4)
	assert.Equal(t, fsEvents[:4], events)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "scp", nil, nil, nil, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "sha256sum", nil, nil, nil, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "", []string{"upload", "download", "rename"}, nil, nil, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "", []string{"rename"}, nil, nil, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "", nil, []string{"SFTP", "HTTP"}, nil, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 4)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "", nil, []string{"SCP"}, nil, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "", nil, nil, []string{"instance1"}, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "", nil, nil, []string{"instance1", "instance3"}, nil, nil, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "", nil, nil, []string{"instance1", "instance2", "instance3"}, nil,
		[]int32{1}, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	data, _, _, err = s.SearchFsEvents(0, 0, "", "", "", nil, nil, []string{"instance1", "instance2", "instance3"}, nil,
		[]int32{1, 2, 3}, 100, 0)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)

	err = sess.Delete(&fsEvents).Error
	assert.NoError(t, err)
}

func TestSearchProviderEvents(t *testing.T) {
	providerEvents := []ProviderEvent{
		{
			ID:         "1",
			Timestamp:  100,
			Action:     "add",
			Username:   "username1",
			IP:         "127.1.1.1",
			ObjectType: "api_key",
			ObjectName: "123",
			ObjectData: []byte("data"),
			InstanceID: "instance1",
		},
		{
			ID:         "2",
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
			ID:         "3",
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
			ID:         "4",
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
			ID:         "5",
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
	_, _, _, err = s.SearchProviderEvents(0, 0, "", "", "", 0, 0, nil, nil, nil, nil)
	assert.ErrorIs(t, err, errNoLimit)
	// test order ASC
	data, sameAtStart, sameAtEnd, err := s.SearchProviderEvents(0, 0, "", "", "", 100, 1, nil, nil, nil, nil)
	assert.NoError(t, err)

	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, providerEvents[0].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 1)
	assert.Equal(t, providerEvents[4].ID, sameAtEnd[0])

	var events []ProviderEvent
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	// test order DESC
	data, sameAtStart, sameAtEnd, err = s.SearchProviderEvents(0, 0, "", "", "", 100, 0, nil, nil, nil, nil)
	assert.NoError(t, err)

	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, providerEvents[4].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 1)
	assert.Equal(t, providerEvents[0].ID, sameAtEnd[0])

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 5)
	// test limit and pagination
	data, sameAtStart, sameAtEnd, err = s.SearchProviderEvents(0, 0, "", "", "", 3, 1, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, providerEvents[0].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 2)
	assert.Equal(t, providerEvents[1].ID, sameAtEnd[1])
	assert.Equal(t, providerEvents[2].ID, sameAtEnd[0])

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
	// get next page
	data, sameAtStart, sameAtEnd, err = s.SearchProviderEvents(events[2].Timestamp, 0, "", "", "", 3, 1, nil, nil,
		nil, sameAtEnd)
	assert.NoError(t, err)
	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, providerEvents[3].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 1)
	assert.Equal(t, providerEvents[4].ID, sameAtEnd[0])

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, providerEvents[3].ID, events[0].ID)
	assert.Equal(t, providerEvents[4].ID, events[1].ID)
	// get previous page
	data, sameAtStart, sameAtEnd, err = s.SearchProviderEvents(0, events[0].Timestamp, "", "", "", 3, 1, nil, nil,
		nil, sameAtStart)
	assert.NoError(t, err)
	assert.Len(t, sameAtStart, 1)
	assert.Equal(t, providerEvents[0].ID, sameAtStart[0])
	assert.Len(t, sameAtEnd, 2)
	assert.Equal(t, providerEvents[1].ID, sameAtEnd[1])
	assert.Equal(t, providerEvents[2].ID, sameAtEnd[0])

	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)
	assert.Equal(t, providerEvents[:3], events)
	// test other search conditions
	data, _, _, err = s.SearchProviderEvents(0, 0, "", "", "", 100, 0, []string{"add", "delete"}, nil,
		nil, nil)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	data, _, _, err = s.SearchProviderEvents(0, 0, "username1", "127.1.1.1", "", 100, 0, []string{"add", "update"}, nil,
		nil, nil)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, providerEvents[0].ID, events[0].ID)

	data, _, _, err = s.SearchProviderEvents(0, 0, "", "", "", 100, 0, nil, []string{"api_key", "user"},
		nil, nil)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 3)

	data, _, _, err = s.SearchProviderEvents(0, 0, "", "", "123", 100, 0, nil, []string{"api_key", "user"},
		nil, nil)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, providerEvents[0].ID, events[0].ID)

	data, _, _, err = s.SearchProviderEvents(0, 0, "", "", "123", 100, 0, nil, []string{"api_key", "user"},
		[]string{"instance2", "instance3"}, nil)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 0)

	data, _, _, err = s.SearchProviderEvents(0, 0, "", "", "", 100, 0, nil, []string{"api_key", "admin"},
		[]string{"instance2", "instance3"}, nil)
	assert.NoError(t, err)
	events = nil
	err = json.Unmarshal(data, &events)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	err = sess.Delete(&providerEvents).Error
	assert.NoError(t, err)
}
