// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"time"

	"github.com/pingcap/tidb/infoschema"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/statistics"
	"github.com/pingcap/tidb/types"
)

// StatsGC is used to GC unnecessary stats.
type StatsGC interface {
	// GCStats will garbage collect the useless stats' info.
	// For dropped tables, we will first update their version
	// so that other tidb could know that table is deleted.
	GCStats(is infoschema.InfoSchema, ddlLease time.Duration) (err error)

	// ClearOutdatedHistoryStats clear outdated historical stats.
	// Only for test.
	ClearOutdatedHistoryStats() error

	// DeleteTableStatsFromKV deletes table statistics from kv.
	// A statsID refers to statistic of a table or a partition.
	DeleteTableStatsFromKV(statsIDs []int64) (err error)
}

// ColStatsTimeInfo records usage information of this column stats.
type ColStatsTimeInfo struct {
	LastUsedAt     *types.Time // last time the column is used
	LastAnalyzedAt *types.Time // last time the column is analyzed
}

// StatsUsage is used to track the usage of column / index statistics.
type StatsUsage interface {
	// Below methods are for predicated columns.

	// LoadColumnStatsUsage returns all columns' usage information.
	LoadColumnStatsUsage(loc *time.Location) (map[model.TableItemID]ColStatsTimeInfo, error)

	// GetPredicateColumns returns IDs of predicate columns, which are the columns whose stats are used(needed) when generating query plans.
	GetPredicateColumns(tableID int64) ([]int64, error)

	// CollectColumnsInExtendedStats returns IDs of the columns involved in extended stats.
	CollectColumnsInExtendedStats(tableID int64) ([]int64, error)

	// Below methods are for index usage.

	// NewSessionIndexUsageCollector creates a new IndexUsageCollector on the list.
	// The returned value's type should be *usage.SessionIndexUsageCollector, use interface{} to avoid cycle import now.
	// TODO: use *usage.SessionIndexUsageCollector instead of interface{}.
	NewSessionIndexUsageCollector() interface{}

	// DumpIndexUsageToKV dumps all collected index usage info to storage.
	DumpIndexUsageToKV() error

	// GCIndexUsage removes unnecessary index usage data.
	GCIndexUsage() error
}

// StatsHistory is used to manage historical stats.
type StatsHistory interface {
	// RecordHistoricalStatsMeta records stats meta of the specified version to stats_meta_history.
	RecordHistoricalStatsMeta(tableID int64, version uint64, source string)

	// CheckHistoricalStatsEnable check whether historical stats is enabled.
	CheckHistoricalStatsEnable() (enable bool, err error)

	// TODO: RecordHistoricalStatsToStorage(dbName string, tableInfo *model.TableInfo, physicalID int64, isPartition bool) (uint64, error)
}

// StatsAnalyze is used to handle auto-analyze and manage analyze jobs.
type StatsAnalyze interface {
	// InsertAnalyzeJob inserts analyze job into mysql.analyze_jobs and gets job ID for further updating job.
	InsertAnalyzeJob(job *statistics.AnalyzeJob, instance string, procID uint64) error

	// DeleteAnalyzeJobs deletes the analyze jobs whose update time is earlier than updateTime.
	DeleteAnalyzeJobs(updateTime time.Time) error

	// TODO: HandleAutoAnalyze
}