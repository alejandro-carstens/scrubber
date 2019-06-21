package responses

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

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
	Failures          []SnapshotShardFailure `json:"failures"`
	Shards            *ShardsInfo            `json:"shards"`
}

func (si *SnapshotInfo) Marshal(container *gabs.Container) (Informable, error) {
	err := json.Unmarshal(container.Bytes(), si)

	return si, err
}

func (si *SnapshotInfo) IsSnapshotInfo() bool {
	return false
}

func (si *SnapshotInfo) Name() string {
	return si.Snapshot
}

func (si *SnapshotInfo) CreationDate() string {
	return si.EndTime
}

type SnapshotShardFailure struct {
	Index     string `json:"index"`
	IndexUUID string `json:"index_uuid"`
	ShardID   int    `json:"shard_id"`
	Reason    string `json:"reason"`
	NodeID    string `json:"node_id"`
	Status    string `json:"status"`
}

type ShardsInfo struct {
	Total      int             `json:"total"`
	Successful int             `json:"successful"`
	Failed     int             `json:"failed"`
	Failures   []*ShardFailure `json:"failures,omitempty"`
}

type ShardFailure struct {
	Index   string                 `json:"_index,omitempty"`
	Shard   int                    `json:"_shard,omitempty"`
	Node    string                 `json:"_node,omitempty"`
	Reason  map[string]interface{} `json:"reason,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Primary bool                   `json:"primary,omitempty"`
}
