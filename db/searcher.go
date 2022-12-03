package db

import (
	"encoding/json"
	"errors"

	"github.com/sftpgo/sdk/plugin/eventsearcher"

	"github.com/sftpgo/sftpgo-plugin-eventsearch/logger"
)

var (
	errNoLimit = errors.New("please specify a limit")
)

type Searcher struct{}

func (s *Searcher) SearchFsEvents(filters *eventsearcher.FsEventSearch) ([]byte, []string, []string, error) {
	if filters.Limit <= 0 {
		return nil, nil, nil, errNoLimit
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	var results []FsEvent
	if filters.StartTimestamp > 0 {
		sess = sess.Where("timestamp >= ?", filters.StartTimestamp)
	}
	if filters.EndTimestamp > 0 {
		sess = sess.Where("timestamp <= ?", filters.EndTimestamp)
	}
	if len(filters.Actions) > 0 {
		sess = sess.Where("action IN ?", filters.Actions)
	}
	if filters.Username != "" {
		sess = sess.Where("username = ?", filters.Username)
	}
	if filters.IP != "" {
		sess = sess.Where("ip = ?", filters.IP)
	}
	if filters.SSHCmd != "" {
		sess = sess.Where("ssh_cmd = ?", filters.SSHCmd)
	}
	if len(filters.Protocols) > 0 {
		sess = sess.Where("protocol IN ?", filters.Protocols)
	}
	if len(filters.InstanceIDs) > 0 {
		sess = sess.Where("instance_id IN ?", filters.InstanceIDs)
	}
	if len(filters.Statuses) > 0 {
		sess = sess.Where("status IN ?", filters.Statuses)
	}
	if len(filters.ExcludeIDs) > 0 {
		sess = sess.Where("id NOT IN ?", filters.ExcludeIDs)
	}
	if filters.FsProvider >= 0 {
		sess = sess.Where("fs_provider = ?", filters.FsProvider)
	}
	if filters.Bucket != "" {
		sess = sess.Where("bucket = ?", filters.Bucket)
	}
	if filters.Endpoint != "" {
		sess = sess.Where("endpoint = ?", filters.Endpoint)
	}
	if filters.Role != "" {
		sess = sess.Where("role = ?", filters.Role)
	}
	sess = sess.Limit(filters.Limit)

	if filters.Order == 0 {
		sess = sess.Order("timestamp DESC, id DESC").Find(&results)
	} else {
		sess = sess.Order("timestamp ASC, id ASC").Find(&results)
	}
	err := sess.Error
	if err != nil {
		logger.AppLogger.Warn("unable to search fs events", "error", err)
		return nil, nil, nil, err
	}

	data, err := json.Marshal(results)
	if err != nil {
		return nil, nil, nil, err
	}

	var sameTsAtStart []string
	var sameTsAtEnd []string

	for idx := range results {
		if results[idx].Timestamp != results[0].Timestamp {
			break
		}
		sameTsAtStart = append(sameTsAtStart, results[idx].ID)
	}
	lastIdx := len(results) - 1
	for i := lastIdx; i >= 0; i-- {
		if results[i].Timestamp != results[lastIdx].Timestamp {
			break
		}
		sameTsAtEnd = append(sameTsAtEnd, results[i].ID)
	}

	return data, sameTsAtStart, sameTsAtEnd, err
}

func (s *Searcher) SearchProviderEvents(filters *eventsearcher.ProviderEventSearch) ([]byte, []string, []string, error) {
	if filters.Limit <= 0 {
		return nil, nil, nil, errNoLimit
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	var results []ProviderEvent
	if filters.StartTimestamp > 0 {
		sess = sess.Where("timestamp >= ?", filters.StartTimestamp)
	}
	if filters.EndTimestamp > 0 {
		sess = sess.Where("timestamp <= ?", filters.EndTimestamp)
	}
	if len(filters.Actions) > 0 {
		sess = sess.Where("action IN ?", filters.Actions)
	}
	if filters.Username != "" {
		sess = sess.Where("username = ?", filters.Username)
	}
	if filters.IP != "" {
		sess = sess.Where("ip = ?", filters.IP)
	}
	if len(filters.ObjectTypes) > 0 {
		sess = sess.Where("object_type IN ?", filters.ObjectTypes)
	}
	if filters.ObjectName != "" {
		sess = sess.Where("object_name = ?", filters.ObjectName)
	}
	if len(filters.InstanceIDs) > 0 {
		sess = sess.Where("instance_id IN ?", filters.InstanceIDs)
	}
	if len(filters.ExcludeIDs) > 0 {
		sess = sess.Where("id NOT IN ?", filters.ExcludeIDs)
	}
	if filters.Role != "" {
		sess = sess.Where("role = ?", filters.Role)
	}
	sess = sess.Limit(filters.Limit)

	if filters.Order == 0 {
		sess = sess.Order("timestamp DESC, id DESC").Find(&results)
	} else {
		sess = sess.Order("timestamp ASC, id ASC").Find(&results)
	}
	err := sess.Error
	if err != nil {
		logger.AppLogger.Warn("unable to search provider events", "error", err)
		return nil, nil, nil, err
	}

	data, err := json.Marshal(results)
	if err != nil {
		return nil, nil, nil, err
	}

	var sameTsAtStart []string
	var sameTsAtEnd []string

	for idx := range results {
		if results[idx].Timestamp != results[0].Timestamp {
			break
		}
		sameTsAtStart = append(sameTsAtStart, results[idx].ID)
	}
	lastIdx := len(results) - 1
	for i := lastIdx; i >= 0; i-- {
		if results[i].Timestamp != results[lastIdx].Timestamp {
			break
		}
		sameTsAtEnd = append(sameTsAtEnd, results[i].ID)
	}

	return data, sameTsAtStart, sameTsAtEnd, err
}
