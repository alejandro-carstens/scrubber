package infos

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
)

// IndexInfo holds an elasticsearch index information
type IndexInfo struct {
	Health                       string `json:"health"`
	Status                       string `json:"status"`
	Index                        string `json:"index"`
	UUID                         string `json:"uuid"`
	Pri                          int    `json:"pri,string"`
	Rep                          int    `json:"rep,string"`
	DocsCount                    int    `json:"docs.count,string"`
	DocsDeleted                  int    `json:"docs.deleted,string"`
	CreationDateInt              int64  `json:"creation.date,string"`
	CreationDateString           string `json:"creation.date.string"`
	StoreSize                    string `json:"store.size"`
	PriStoreSize                 string `json:"pri.store.size"`
	CompletionSize               string `json:"completion.size"`
	PriCompletionSize            string `json:"pri.completion.size"`
	FielddataMemorySize          string `json:"fielddata.memory_size"`
	PriFielddataMemorySize       string `json:"pri.fielddata.memory_size"`
	FielddataEvictions           int    `json:"fielddata.evictions,string"`
	PriFielddataEvictions        int    `json:"pri.fielddata.evictions,string"`
	QueryCacheMemorySize         string `json:"query_cache.memory_size"`
	PriQueryCacheMemorySize      string `json:"pri.query_cache.memory_size"`
	QueryCacheEvictions          int    `json:"query_cache.evictions,string"`
	PriQueryCacheEvictions       int    `json:"pri.query_cache.evictions,string"`
	RequestCacheMemorySize       string `json:"request_cache.memory_size"`
	PriRequestCacheMemorySize    string `json:"pri.request_cache.memory_size"`
	RequestCacheEvictions        int    `json:"request_cache.evictions,string"`
	PriRequestCacheEvictions     int    `json:"pri.request_cache.evictions,string"`
	RequestCacheHitCount         int    `json:"request_cache.hit_count,string"`
	PriRequestCacheHitCount      int    `json:"pri.request_cache.hit_count,string"`
	RequestCacheMissCount        int    `json:"request_cache.miss_count,string"`
	PriRequestCacheMissCount     int    `json:"pri.request_cache.miss_count,string"`
	FlushTotal                   int    `json:"flush.total"`
	PriFlushTotal                int    `json:"pri.flush.total"`
	FlushTotalTime               string `json:"flush.total_time"`
	PriFlushTotalTime            string `json:"pri.flush.total_time"`
	GetCurrent                   int    `json:"get.current,string"`
	PriGetCurrent                int    `json:"pri.get.current,string"`
	GetTime                      string `json:"get.time"`
	PriGetTime                   string `json:"pri.get.time"`
	GetTotal                     int    `json:"get.total,string"`
	PriGetTotal                  int    `json:"pri.get.total,string"`
	GetExistsTime                string `json:"get.exists_time"`
	PriGetExistsTime             string `json:"pri.get.exists_time"`
	GetExistsTotal               int    `json:"get.exists_total,string"`
	PriGetExistsTotal            int    `json:"pri.get.exists_total,string"`
	GetMissingTime               string `json:"get.missing_time"`
	PriGetMissingTime            string `json:"pri.get.missing_time"`
	GetMissingTotal              int    `json:"get.missing_total,string"`
	PriGetMissingTotal           int    `json:"pri.get.missing_total,string"`
	IndexingDeleteCurrent        int    `json:"indexing.delete_current,string"`
	PriIndexingDeleteCurrent     int    `json:"pri.indexing.delete_current,string"`
	IndexingDeleteTime           string `json:"indexing.delete_time"`
	PriIndexingDeleteTime        string `json:"pri.indexing.delete_time"`
	IndexingDeleteTotal          int    `json:"indexing.delete_total,string"`
	PriIndexingDeleteTotal       int    `json:"pri.indexing.delete_total,string"`
	IndexingIndexCurrent         int    `json:"indexing.index_current,string"`
	PriIndexingIndexCurrent      int    `json:"pri.indexing.index_current,string"`
	IndexingIndexTime            string `json:"indexing.index_time"`
	PriIndexingIndexTime         string `json:"pri.indexing.index_time"`
	IndexingIndexTotal           int    `json:"indexing.index_total,string"`
	PriIndexingIndexTotal        int    `json:"pri.indexing.index_total,string"`
	IndexingIndexFailed          int    `json:"indexing.index_failed,string"`
	PriIndexingIndexFailed       int    `json:"pri.indexing.index_failed,string"`
	MergesCurrent                int    `json:"merges.current,string"`
	PriMergesCurrent             int    `json:"pri.merges.current,string"`
	MergesCurrentDocs            int    `json:"merges.current_docs,string"`
	PriMergesCurrentDocs         int    `json:"pri.merges.current_docs,string"`
	MergesCurrentSize            string `json:"merges.current_size"`
	PriMergesCurrentSize         string `json:"pri.merges.current_size"`
	MergesTotal                  int    `json:"merges.total,string"`
	PriMergesTotal               int    `json:"pri.merges.total,string"`
	MergesTotalDocs              int    `json:"merges.total_docs,string"`
	PriMergesTotalDocs           int    `json:"pri.merges.total_docs,string"`
	MergesTotalSize              string `json:"merges.total_size"`
	PriMergesTotalSize           string `json:"pri.merges.total_size"`
	MergesTotalTime              string `json:"merges.total_time"`
	PriMergesTotalTime           string `json:"pri.merges.total_time"`
	RefreshTotal                 int    `json:"refresh.total,string"`
	PriRefreshTotal              int    `json:"pri.refresh.total,string"`
	RefreshTime                  string `json:"refresh.time"`
	PriRefreshTime               string `json:"pri.refresh.time"`
	RefreshListeners             int    `json:"refresh.listeners,string"`
	PriRefreshListeners          int    `json:"pri.refresh.listeners,string"`
	SearchFetchCurrent           int    `json:"search.fetch_current,string"`
	PriSearchFetchCurrent        int    `json:"pri.search.fetch_current,string"`
	SearchFetchTime              string `json:"search.fetch_time"`
	PriSearchFetchTime           string `json:"pri.search.fetch_time"`
	SearchFetchTotal             int    `json:"search.fetch_total,string"`
	PriSearchFetchTotal          int    `json:"pri.search.fetch_total,string"`
	SearchOpenContexts           int    `json:"search.open_contexts,string"`
	PriSearchOpenContexts        int    `json:"pri.search.open_contexts,string"`
	SearchQueryCurrent           int    `json:"search.query_current,string"`
	PriSearchQueryCurrent        int    `json:"pri.search.query_current,string"`
	SearchQueryTime              string `json:"search.query_time"`
	PriSearchQueryTime           string `json:"pri.search.query_time"`
	SearchQueryTotal             int    `json:"search.query_total,string"`
	PriSearchQueryTotal          int    `json:"pri.search.query_total,string"`
	SearchScrollCurrent          int    `json:"search.scroll_current,string"`
	PriSearchScrollCurrent       int    `json:"pri.search.scroll_current,string"`
	SearchScrollTime             string `json:"search.scroll_time"`
	PriSearchScrollTime          string `json:"pri.search.scroll_time"`
	SearchScrollTotal            int    `json:"search.scroll_total,string"`
	PriSearchScrollTotal         int    `json:"pri.search.scroll_total,string"`
	SegmentsCount                int    `json:"segments.count,string"`
	PriSegmentsCount             int    `json:"pri.segments.count,string"`
	SegmentsMemory               string `json:"segments.memory"`
	PriSegmentsMemory            string `json:"pri.segments.memory"`
	SegmentsIndexWriterMemory    string `json:"segments.index_writer_memory"`
	PriSegmentsIndexWriterMemory string `json:"pri.segments.index_writer_memory"`
	SegmentsVersionMapMemory     string `json:"segments.version_map_memory"`
	PriSegmentsVersionMapMemory  string `json:"pri.segments.version_map_memory"`
	SegmentsFixedBitsetMemory    string `json:"segments.fixed_bitset_memory"`
	PriSegmentsFixedBitsetMemory string `json:"pri.segments.fixed_bitset_memory"`
	WarmerCurrent                int    `json:"warmer.count,string"`
	PriWarmerCurrent             int    `json:"pri.warmer.count,string"`
	WarmerTotal                  int    `json:"warmer.total,string"`
	PriWarmerTotal               int    `json:"pri.warmer.total,string"`
	WarmerTotalTime              string `json:"warmer.total_time"`
	PriWarmerTotalTime           string `json:"pri.warmer.total_time"`
	SuggestCurrent               int    `json:"suggest.current,string"`
	PriSuggestCurrent            int    `json:"pri.suggest.current,string"`
	SuggestTime                  string `json:"suggest.time"`
	PriSuggestTime               string `json:"pri.suggest.time"`
	SuggestTotal                 int    `json:"suggest.total,string"`
	PriSuggestTotal              int    `json:"pri.suggest.total,string"`
	MemoryTotal                  string `json:"memory.total"`
	PriMemoryTotal               string `json:"pri.memory.total"`
}

// Marshal implementation of the Informable interface
func (ii *IndexInfo) Marshal(container *gabs.Container) (Informable, error) {
	err := json.Unmarshal(container.Bytes(), ii)

	return ii, err
}

// IsSnapshotInfo implementation of the Informable interface
func (ii *IndexInfo) IsSnapshotInfo() bool {
	return false
}

// Name implementation of the Informable interface
func (ii *IndexInfo) Name() string {
	return ii.Index
}

// CreationDate implementation of the Informable interface
func (ii *IndexInfo) CreationDate() string {
	return ii.CreationDateString
}
