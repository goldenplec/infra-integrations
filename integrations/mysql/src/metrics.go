package main

import "github.com/newrelic/infra-integrations-sdk/metric"

var defaultMetrics = map[string][]interface{}{
	"provider.abortedClientsPerSecond":                 {"Aborted_clients", metric.COUNTER},
	"provider.abortedConnectsPerSecond":                {"Aborted_connects", metric.COUNTER},
	"provider.bytesReceivedPerSecond":                  {"Bytes_received", metric.COUNTER},
	"provider.bytesSentPerSecond":                      {"Bytes_sent", metric.COUNTER},
	"provider.comDeletePerSecond":                      {"Com_delete", metric.COUNTER},
	"provider.comDeleteMultiPerSecond":                 {"Com_delete_multi", metric.COUNTER},
	"provider.comInsertPerSecond":                      {"Com_insert", metric.COUNTER},
	"provider.comInsertSelectPerSecond":                {"Com_insert_select", metric.COUNTER},
	"provider.comReplaceSelectPerSecond":               {"Com_replace_select", metric.COUNTER},
	"provider.comSelectPerSecond":                      {"Com_select", metric.COUNTER},
	"provider.comUpdatePerSecond":                      {"Com_update", metric.COUNTER},
	"provider.comUpdateMultiPerSecond":                 {"Com_update_multi", metric.COUNTER},
	"provider.connectionErrorsMaxConnectionsPerSecond": {"Connection_errors_max_connections", metric.COUNTER},
	"provider.connectionsPerSecond":                    {"Connections", metric.COUNTER},
	"provider.handlerRollbackPerSecond":                {"Handler_rollback", metric.COUNTER},
	"provider.innodbBufferPoolPagesData":               {"Innodb_buffer_pool_pages_data", metric.GAUGE},
	"provider.innodbBufferPoolPagesFree":               {"Innodb_buffer_pool_pages_free", metric.GAUGE},
	"provider.innodbBufferPoolPagesTotal":              {"Innodb_buffer_pool_pages_total", metric.GAUGE},
	"provider.innodbDataReadPerSecond":                 {"Innodb_data_read", metric.COUNTER},
	"provider.innodbDataWrittenPerSecond":              {"Innodb_data_written", metric.COUNTER},
	"provider.innodbLogWaitsPerSecond":                 {"Innodb_log_waits", metric.COUNTER},
	"provider.innodbRowLockCurrentWaits":               {"Innodb_row_lock_current_waits", metric.GAUGE},
	"provider.innodbRowLockTimeAvg":                    {"Innodb_row_lock_time_avg", metric.GAUGE},
	"provider.innodbRowLockWaitsPerSecond":             {"Innodb_row_lock_waits", metric.COUNTER},
	"provider.maxConnections":                          {"Max_used_connections", metric.GAUGE},
	"provider.openFiles":                               {"Open_files", metric.GAUGE},
	"provider.openTables":                              {"Open_tables", metric.GAUGE},
	"provider.openedTablesPerSecond":                   {"Opened_tables", metric.COUNTER},
	"provider.preparedStmtCountPerSecond":              {"Prepared_stmt_count", metric.COUNTER},
	"provider.qCacheFreeMemory":                        {"Qcache_free_memory", metric.GAUGE},
	"provider.qCacheNotCachedPerSecond":                {"Qcache_not_cached", metric.COUNTER},
	"provider.queriesPerSecond":                        {"Queries", metric.COUNTER},
	"provider.questionsPerSecond":                      {"Questions", metric.COUNTER},
	"provider.slowQueriesPerSecond":                    {"Slow_queries", metric.COUNTER},
	"provider.tablesLocksWaitedPerSecond":              {"Table_locks_waited", metric.COUNTER},
	"provider.threadsConnected":                        {"Threads_connected", metric.GAUGE},
	"provider.threadsRunning":                          {"Threads_running", metric.GAUGE},
	"provider.qCacheUtilization":                       {qCacheUtilization, metric.GAUGE},
	"provider.qCacheHitRatio":                          {qCacheHitRatio, metric.GAUGE},
	"software.Edition":                                 {"version_comment", metric.ATTRIBUTE},
	"software.Version":                                 {"version", metric.ATTRIBUTE},
	"cluster.NodeType":                                 {"node_type", metric.ATTRIBUTE},
}

func qCacheUtilization(metrics map[string]interface{}) (float64, bool) {
	qCacheFreeBlocks, ok1 := metrics["Qcache_free_blocks"].(int)
	qCacheTotalBlocks, ok2 := metrics["Qcache_total_blocks"].(int)

	if ok1 && ok2 {
		return 1 - (float64(qCacheFreeBlocks) / float64(qCacheTotalBlocks)), true
	}
	return 0, false
}

func qCacheHitRatio(metrics map[string]interface{}) (float64, bool) {
	qCacheHits, ok1 := metrics["Qcache_hits"].(int)
	queries, ok2 := metrics["Queries"].(int)

	if ok1 && ok2 {
		return float64(qCacheHits) / float64(queries), true
	}
	return 0, false
}

var extendedMetrics = map[string][]interface{}{
	"provider.createdTmpDiskTablesPerSecond":     {"Created_tmp_disk_tables", metric.COUNTER},
	"provider.createdTmpFilesPerSecond":          {"Created_tmp_files", metric.COUNTER},
	"provider.createdTmpTablesPerSecond":         {"Created_tmp_tables", metric.COUNTER},
	"provider.handlerDeletePerSecond":            {"Handler_delete", metric.COUNTER},
	"provider.handlerReadFirstPerSecond":         {"Handler_read_first", metric.COUNTER},
	"provider.handlerReadKeyPerSecond":           {"Handler_read_key", metric.COUNTER},
	"provider.handlerReadRndPerSecond":           {"Handler_read_rnd", metric.COUNTER},
	"provider.handlerReadRndNextPerSecond":       {"Handler_read_rnd_next", metric.COUNTER},
	"provider.handlerUpdatePerSecond":            {"Handler_update", metric.COUNTER},
	"provider.handlerWritePerSecond":             {"Handler_write", metric.COUNTER},
	"provider.maxExecutuibTimeExceededPerSecond": {"Max_execution_time_exceeded", metric.COUNTER},
	"provider.qCacheFreeBlocks":                  {"Qcache_free_blocks", metric.GAUGE},
	"provider.qCacheHitsPerSecond":               {"Qcache_hits", metric.COUNTER},
	"provider.qCacheInserts":                     {"Qcache_inserts", metric.GAUGE},
	"provider.qCacheLowmemPrunesPerSecond":       {"Qcache_lowmem_prunes", metric.COUNTER},
	"provider.qCacheQueriesInCachePerSecond":     {"Qcache_queries_in_cache", metric.COUNTER},
	"provider.qCacheTotalBlocks":                 {"Qcache_total_blocks", metric.GAUGE},
	"provider.selectFullJoinPerSecond":           {"Select_full_join", metric.COUNTER},
	"provider.selectFullJoinRangePerSecond":      {"Select_full_range_join", metric.COUNTER},
	"provider.selectRangePerSecond":              {"Select_range", metric.COUNTER},
	"provider.selectRangeCheckPerSecond":         {"Select_range_check", metric.COUNTER},
	"provider.sortMergePassesPerSecond":          {"Sort_merge_passes", metric.COUNTER},
	"provider.sortRangePerSecond":                {"Sort_range", metric.COUNTER},
	"provider.sortRowsPerSecond":                 {"Sort_rows", metric.COUNTER},
	"provider.sortScanPerSecond":                 {"Sort_scan", metric.COUNTER},
	"provider.tableOpenCacheHitsPerSecond":       {"Table_open_cache_hits", metric.COUNTER},
	"provider.tableOpenCacheMissesPerSecond":     {"Table_open_cache_misses", metric.COUNTER},
	"provider.tableOpenCacheOverflowsPerSecond":  {"Table_open_cache_overflows", metric.COUNTER},
	"provider.threadsCached":                     {"Threads_cached", metric.GAUGE},
	"provider.threadsCreatedPerSecond":           {"Threads_created", metric.COUNTER},
	"provider.threadCacheMissRate":               {threadCacheMissRate, metric.GAUGE},
}

func threadCacheMissRate(metrics map[string]interface{}) (float64, bool) {
	threadsCreated, ok1 := metrics["Threads_created"].(int)
	connections, ok2 := metrics["Connections"].(int)

	if ok1 && ok2 {
		return float64(threadsCreated) / float64(connections), true
	}
	return 0, false
}

var innodbMetrics = map[string][]interface{}{
	"provider.innodbBufferPoolPagesDirty":                {"Innodb_buffer_pool_pages_dirty", metric.GAUGE},
	"provider.innodbBufferPoolPagesFlushedPerSecond":     {"Innodb_buffer_pool_pages_flushed", metric.COUNTER},
	"provider.innodbBufferPoolReadAheadPerSecond":        {"Innodb_buffer_pool_read_ahead", metric.COUNTER},
	"provider.innodbBufferPoolReadAheadEvictedPerSecond": {"Innodb_buffer_pool_read_ahead_evicted", metric.COUNTER},
	"provider.innodbBufferPoolReadAheadRndPerSecond":     {"Innodb_buffer_pool_read_ahead_rnd", metric.COUNTER},
	"provider.innodbBufferPoolReadRequestsPerSecond":     {"Innodb_buffer_pool_read_requests", metric.COUNTER},
	"provider.innodbBufferPoolReads":                     {"Innodb_buffer_pool_reads", metric.GAUGE},
	"provider.innodbBufferPoolWaitFreePerSecond":         {"Innodb_buffer_pool_wait_free", metric.COUNTER},
	"provider.innodbBufferPoolWriteRequestsPerSecond":    {"Innodb_buffer_pool_write_requests", metric.COUNTER},
	"provider.innodbDataFsyncsPerSecond":                 {"Innodb_data_fsyncs", metric.COUNTER},
	"provider.innodbDataPendingFsyncsPerSecond":          {"Innodb_data_pending_fsyncs", metric.COUNTER},
	"provider.innodbDataPendingReadsPerSecond":           {"Innodb_data_pending_reads", metric.COUNTER},
	"provider.innodbDataPendingWritesPerSecond":          {"Innodb_data_pending_writes", metric.COUNTER},
	"provider.innodbDataReadsPerSecond":                  {"Innodb_data_reads", metric.COUNTER},
	"provider.innodbDataWritesPerSecond":                 {"Innodb_data_writes", metric.COUNTER},
	"provider.innodbLogWriteRequestsPerSecond":           {"Innodb_log_write_requests", metric.COUNTER},
	"provider.innodbWritesPerSecond":                     {"Innodb_log_writes", metric.COUNTER},
	"provider.innodbNumOpenFiles":                        {"Innodb_num_open_files", metric.GAUGE},
	"provider.innodbOsLogFsyncsPerSecond":                {"Innodb_os_log_fsyncs", metric.COUNTER},
	"provider.innodbOsLogPendingFsyncsPerSecond":         {"Innodb_os_log_pending_fsyncs", metric.COUNTER},
	"provider.innodbOsLogPendingWritesPerSecond":         {"Innodb_os_log_pending_writes", metric.COUNTER},
	"provider.innodbOsLogWrittenPerSecond":               {"Innodb_os_log_written", metric.COUNTER},
	"provider.innodbPagesCreatedPerSecond":               {"Innodb_pages_created", metric.COUNTER},
	"provider.innodbPagesReadPerSecond":                  {"Innodb_pages_read", metric.COUNTER},
	"provider.innodbPagesWrittenPerSecond":               {"Innodb_pages_written", metric.COUNTER},
	"provider.innodbRowsDeletedPerSecond":                {"Innodb_rows_deleted", metric.COUNTER},
	"provider.innodbRowsInsertedPerSecond":               {"Innodb_rows_inserted", metric.COUNTER},
	"provider.innodbRowsReadPerSecond":                   {"Innodb_rows_read", metric.COUNTER},
	"provider.innodbRowsUpdatedPerSecond":                {"Innodb_rows_updated", metric.COUNTER},
}

var myisamMetrics = map[string][]interface{}{
	"provider.keyBlocksNotFlushed":       {"Key_blocks_not_flushed", metric.GAUGE},
	"provider.keyCacheUtilization":       {keyCacheUtilization, metric.GAUGE},
	"provider.keyReadRequestsPerSecond":  {"Key_read_requests", metric.COUNTER},
	"provider.KeyReadsPerSecond":         {"Key_reads", metric.COUNTER},
	"provider.KeyWriteRequestsPerSecond": {"Key_write_requests", metric.COUNTER},
	"provider.KeyWritesPerSecond":        {"Key_writes", metric.COUNTER},
}

func keyCacheUtilization(metrics map[string]interface{}) (float64, bool) {
	keyBlocksUnused, ok1 := metrics["Key_blocks_unused"].(int)
	keyCacheBlockSize, ok2 := metrics["key_cache_block_size"].(int)
	keyBufferSize, ok3 := metrics["key_buffer_size"].(int)

	if ok1 && ok2 && ok3 {
		return 1 - (float64(keyBlocksUnused) * float64(keyCacheBlockSize) / float64(keyBufferSize)), true
	}
	return 0, false
}
