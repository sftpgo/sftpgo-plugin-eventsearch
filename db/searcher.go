package db

import (
	"encoding/json"
	"errors"

	"github.com/sftpgo/sftpgo-plugin-eventsearch/logger"
)

var (
	errNoLimit = errors.New("please specify a limit")
)

type Searcher struct{}

func (s *Searcher) SearchFsEvents(startTimestamp, endTimestamp int64, username, ip, sshCmd string, actions,
	protocols, instanceIDs, excludeIDs []string, statuses []int32, limit, order int,
) ([]byte, []string, []string, error) {
	if limit <= 0 {
		return nil, nil, nil, errNoLimit
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	var results []FsEvent
	if startTimestamp > 0 {
		sess = sess.Where("timestamp >= ?", startTimestamp)
	}
	if endTimestamp > 0 {
		sess = sess.Where("timestamp <= ?", endTimestamp)
	}
	if len(actions) > 0 {
		sess = sess.Where("action IN ?", actions)
	}
	if username != "" {
		sess = sess.Where("username = ?", username)
	}
	if ip != "" {
		sess = sess.Where("ip = ?", ip)
	}
	if sshCmd != "" {
		sess = sess.Where("ssh_cmd = ?", sshCmd)
	}
	if len(protocols) > 0 {
		sess = sess.Where("protocol IN ?", protocols)
	}
	if len(instanceIDs) > 0 {
		sess = sess.Where("instance_id IN ?", instanceIDs)
	}
	if len(statuses) > 0 {
		sess = sess.Where("status IN ?", statuses)
	}
	if len(excludeIDs) > 0 {
		sess = sess.Where("id NOT IN ?", excludeIDs)
	}
	sess = sess.Limit(limit)

	if order == 0 {
		sess = sess.Order("timestamp DESC").Find(&results)
	} else {
		sess = sess.Order("timestamp ASC").Find(&results)
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

func (s *Searcher) SearchProviderEvents(startTimestamp, endTimestamp int64, username, ip, objectName string,
	limit, order int, actions, objectTypes, instanceIDs, excludeIDs []string,
) ([]byte, []string, []string, error) {
	if limit <= 0 {
		return nil, nil, nil, errNoLimit
	}

	sess, cancel := getDefaultSession()
	defer cancel()

	var results []ProviderEvent
	if startTimestamp > 0 {
		sess = sess.Where("timestamp >= ?", startTimestamp)
	}
	if endTimestamp > 0 {
		sess = sess.Where("timestamp <= ?", endTimestamp)
	}
	if len(actions) > 0 {
		sess = sess.Where("action IN ?", actions)
	}
	if username != "" {
		sess = sess.Where("username = ?", username)
	}
	if ip != "" {
		sess = sess.Where("ip = ?", ip)
	}
	if len(objectTypes) > 0 {
		sess = sess.Where("object_type IN ?", objectTypes)
	}
	if objectName != "" {
		sess = sess.Where("object_name = ?", objectName)
	}
	if len(instanceIDs) > 0 {
		sess = sess.Where("instance_id IN ?", instanceIDs)
	}
	if len(excludeIDs) > 0 {
		sess = sess.Where("id NOT IN ?", excludeIDs)
	}
	sess = sess.Limit(limit)

	if order == 0 {
		sess = sess.Order("timestamp DESC").Find(&results)
	} else {
		sess = sess.Order("timestamp ASC").Find(&results)
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
