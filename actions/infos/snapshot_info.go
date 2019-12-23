package infos

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

// SnapshotInfo holds an elasticsearch snapshot information
type SnapshotInfo struct {
	Snapshot          string                 `json:"snapshot"`
	UUID              string                 `json:"uuid"`
	VersionID         int                    `json:"version_id"`
	Version           string                 `json:"version"`
	Indices           []string               `json:"indices"`
	State             string                 `json:"state"`
	Reason            string                 `json:"reason"`
	StartTime         string                 `json:"start_time"`
	StartTimeInMillis int64                  `json:"start_time_in_millis"`
	EndTime           string                 `json:"end_time"`
	EndTimeInMillis   int64                  `json:"end_time_in_millis"`
	DurationInMillis  int64                  `json:"duration_in_millis"`
	Failures          []snapshotShardFailure `json:"failures"`
	Shards            *shardsInfo            `json:"shards"`
}

// Marshal implementation of the Informable interface
func (si *SnapshotInfo) Marshal(container *gabs.Container) (Informable, error) {
	err := json.Unmarshal(container.Bytes(), si)

	return si, err
}

// IsSnapshot implementation of the Informable interface
func (si *SnapshotInfo) IsSnapshotInfo() bool {
	return true
}

// Name implementation of the Informable interface
func (si *SnapshotInfo) Name() string {
	return si.Snapshot
}

// CreationDate implementation of the Informable interface
func (si *SnapshotInfo) CreationDate() string {
	return si.EndTime
}

type snapshotShardFailure struct {
	Index     string `json:"index"`
	IndexUUID string `json:"index_uuid"`
	ShardID   int    `json:"shard_id"`
	Reason    string `json:"reason"`
	NodeID    string `json:"node_id"`
	Status    string `json:"status"`
}

type shardsInfo struct {
	Total      int             `json:"total"`
	Successful int             `json:"successful"`
	Failed     int             `json:"failed"`
	Failures   []*shardFailure `json:"failures,omitempty"`
}

type shardFailure struct {
	Index   string                 `json:"_index,omitempty"`
	Shard   int                    `json:"_shard,omitempty"`
	Node    string                 `json:"_node,omitempty"`
	Reason  map[string]interface{} `json:"reason,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Primary bool                   `json:"primary,omitempty"`
}
