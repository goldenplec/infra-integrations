package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/jmx"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
)

// All metrics we want to provide for the cassandra integration
var metricsDefinition = map[string][]interface{}{
	"provider.viewWriteLatencyPerSecond":                    {"org.apache.cassandra.metrics:type=ClientRequest,scope=ViewWrite,name=Latency,attr=OneMinuteRate", metric.GAUGE},
	"provider.rangeSliceLatencyPerSecond":                   {"org.apache.cassandra.metrics:type=ClientRequest,scope=RangeSlice,name=Latency,attr=OneMinuteRate", metric.GAUGE},
	"provider.CASWriteLatencyPerSecond":                     {"org.apache.cassandra.metrics:type=ClientRequest,scope=CASWrite,name=Latency,attr=OneMinuteRate", metric.GAUGE},
	"provider.readLatencyPerSecond":                         {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Latency,attr=OneMinuteRate", metric.GAUGE},
	"provider.CASReadLatencyPerSecond":                      {"org.apache.cassandra.metrics:type=ClientRequest,scope=CASRead,name=Latency,attr=OneMinuteRate", metric.GAUGE},
	"provider.writeLatencyPerSecond":                        {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Latency,attr=OneMinuteRate", metric.GAUGE},
	"provider.writeLatency98thPercentile":                   {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Latency,attr=98thPercentile", metric.GAUGE},
	"provider.writeLatency99thPercentile":                   {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Latency,attr=99thPercentile", metric.GAUGE},
	"provider.writeLatency999thPercentile":                  {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Latency,attr=999thPercentile", metric.GAUGE},
	"provider.writeLatency50thPercentile":                   {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Latency,attr=50thPercentile", metric.GAUGE},
	"provider.writeLatency75thPercentile":                   {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Latency,attr=75thPercentile", metric.GAUGE},
	"provider.writeLatency95thPercentile":                   {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Latency,attr=95thPercentile", metric.GAUGE},
	"provider.readLatency98thPercentile":                    {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Latency,attr=98thPercentile", metric.GAUGE},
	"provider.readLatency99thPercentile":                    {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Latency,attr=99thPercentile", metric.GAUGE},
	"provider.readLatency999thPercentile":                   {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Latency,attr=999thPercentile", metric.GAUGE},
	"provider.readLatency50thPercentile":                    {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Latency,attr=50thPercentile", metric.GAUGE},
	"provider.readLatency75thPercentile":                    {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Latency,attr=75thPercentile", metric.GAUGE},
	"provider.readLatency95thPercentile":                    {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Latency,attr=95thPercentile", metric.GAUGE},
	"provider.requestCounterMutationStageActiveTasks":       {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=CounterMutationStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.requestCounterMutationStagePendingTasks":      {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=CounterMutationStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.requestMutationStageActiveTasks":              {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=MutationStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.requestMutationStagePendingTasks":             {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=MutationStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.requestReadRepairStageActiveTasks":            {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=ReadRepairStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.requestReadRepairStagePendingTasks":           {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=ReadRepairStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.requestReadStageActiveTasks":                  {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=ReadStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.requestReadStagePendingTasks":                 {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=ReadStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.requestRequestResponseStageActiveTasks":       {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=RequestResponseStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.requestRequestResponseStagePendingTasks":      {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=RequestResponseStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.requestViewMutationStageActiveTasks":          {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=ViewMutationStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.requestViewMutationStagePendingTasks":         {"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=ViewMutationStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalAntiEntropyStageActiveTasks":          {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=AntiEntropyStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalAntiEntropyStagePendingTasks":         {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=AntiEntropyStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalCacheCleanupExecutorActiveTasks":      {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=CacheCleanupExecutor,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalCacheCleanupExecutorPendingTasks":     {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=CacheCleanupExecutor,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalCompactionExecutorActiveTasks":        {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=CompactionExecutor,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalCompactionExecutorPendingTasks":       {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=CompactionExecutor,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalGossipStageActiveTasks":               {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=GossipStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalGossipStagePendingTasks":              {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=GossipStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalHintsDispatcherActiveTasks":           {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=HintsDispatcher,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalHintsDispatcherPendingTasks":          {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=HintsDispatcher,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalInternalResponseStageActiveTasks":     {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=InternalResponseStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalInternalResponseStagePendingTasks":    {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=InternalResponseStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalMemtableFlushWriterActiveTasks":       {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MemtableFlushWriter,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalMemtableFlushWriterPendingTasks":      {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MemtableFlushWriter,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalMemtablePostFlushActiveTasks":         {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MemtablePostFlush,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalMemtablePostFlushPendingTasks":        {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MemtablePostFlush,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalMemtableReclaimMemoryActiveTasks":     {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MemtableReclaimMemory,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalMemtableReclaimMemoryPendingTasks":    {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MemtableReclaimMemory,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalMigrationStageActiveTasks":            {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MigrationStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalMigrationStagePendingTasks":           {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MigrationStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalMiscStageActiveTasks":                 {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MiscStage,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalMiscStagePendingTasks":                {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=MiscStage,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalPendingRangeCalculatorActiveTasks":    {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=PendingRangeCalculator,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalPendingRangeCalculatorPendingTasks":   {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=PendingRangeCalculator,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalSamplerActiveTasks":                   {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=Sampler,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalSamplerPendingTasks":                  {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=Sampler,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalSecondaryIndexManagementActiveTasks":  {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=SecondaryIndexManagement,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalSecondaryIndexManagementPendingTasks": {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=SecondaryIndexManagement,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.internalValidationExecutorActiveTasks":        {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=ValidationExecutor,name=ActiveTasks,attr=Value", metric.GAUGE},
	"provider.internalValidationExecutorPendingTasks":       {"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=ValidationExecutor,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.droppedBatchRemoveMessagesCount":              {"org.apache.cassandra.metrics:type=DroppedMessage,scope=BATCH_REMOVE,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedBatchStoreMessagesCount":               {"org.apache.cassandra.metrics:type=DroppedMessage,scope=BATCH_STORE,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedCounterMutationMessagesCount":          {"org.apache.cassandra.metrics:type=DroppedMessage,scope=COUNTER_MUTATION,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedHintMessagesCount":                     {"org.apache.cassandra.metrics:type=DroppedMessage,scope=HINT,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedMutationMessagesCount":                 {"org.apache.cassandra.metrics:type=DroppedMessage,scope=MUTATION,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedPagedRangeMessagesCount":               {"org.apache.cassandra.metrics:type=DroppedMessage,scope=PAGED_RANGE,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedRangeSliceMessagesCount":               {"org.apache.cassandra.metrics:type=DroppedMessage,scope=RANGE_SLICE,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedReadMessagesCount":                     {"org.apache.cassandra.metrics:type=DroppedMessage,scope=READ,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedReadRepairMessagesCount":               {"org.apache.cassandra.metrics:type=DroppedMessage,scope=READ_REPAIR,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedRequestResponseMessagesCount":          {"org.apache.cassandra.metrics:type=DroppedMessage,scope=REQUEST_RESPONSE,name=Dropped,attr=Count", metric.GAUGE},
	"provider.droppedTraceMessagesCount":                    {"org.apache.cassandra.metrics:type=DroppedMessage,scope=_TRACE,name=Dropped,attr=Count", metric.GAUGE},
	"provider.liveSSTableCount":                             {"org.apache.cassandra.metrics:type=ColumnFamily,name=LiveSSTableCount,attr=Value", metric.GAUGE},
	"provider.totalHints":                                   {"org.apache.cassandra.metrics:type=Storage,name=TotalHints,attr=Count", metric.GAUGE},
	"provider.totalHintsInProgress":                         {"org.apache.cassandra.metrics:type=Storage,name=TotalHintsInProgress,attr=Count", metric.GAUGE},

	"provider.keyCacheCapacityInBytes":          {"org.apache.cassandra.metrics:type=Cache,scope=KeyCache,name=Capacity,attr=Value", metric.GAUGE},
	"provider.keyCacheHitsPerSecond":            {"org.apache.cassandra.metrics:type=Cache,scope=KeyCache,name=Hits,attr=OneMinuteRate", metric.GAUGE},
	"provider.keyCacheHitRatio":                 {"org.apache.cassandra.metrics:type=Cache,scope=KeyCache,name=HitRate,attr=Value", metric.GAUGE},
	"provider.keyCacheRequestsPerSecond":        {"org.apache.cassandra.metrics:type=Cache,scope=KeyCache,name=Requests,attr=OneMinuteRate", metric.GAUGE},
	"provider.keyCacheSize":                     {"org.apache.cassandra.metrics:type=Cache,scope=KeyCache,name=Size,attr=Value", metric.GAUGE},
	"provider.rowCacheCapacityInBytes":          {"org.apache.cassandra.metrics:type=Cache,scope=RowCache,name=Capacity,attr=Value", metric.GAUGE},
	"provider.rowCacheHits":                     {"org.apache.cassandra.metrics:type=Cache,scope=RowCache,name=Hits,attr=OneMinuteRate", metric.GAUGE},
	"provider.rowCacheHitRatio":                 {"org.apache.cassandra.metrics:type=Cache,scope=RowCache,name=HitRate,attr=Value", metric.GAUGE},
	"provider.rowCacheRequests":                 {"org.apache.cassandra.metrics:type=Cache,scope=RowCache,name=Requests,attr=OneMinuteRate", metric.GAUGE},
	"provider.rowCacheSize":                     {"org.apache.cassandra.metrics:type=Cache,scope=RowCache,name=Size,attr=Value", metric.GAUGE},
	"provider.readTimeoutsPerSecond":            {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Timeouts,attr=OneMinuteRate", metric.GAUGE},
	"provider.readUnavailablesPerSecond":        {"org.apache.cassandra.metrics:type=ClientRequest,scope=Read,name=Unavailables,attr=OneMinuteRate", metric.GAUGE},
	"provider.writeTimeoutsPerSecond":           {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Timeouts,attr=OneMinuteRate", metric.COUNTER},
	"provider.writeUnavailablesPerSecond":       {"org.apache.cassandra.metrics:type=ClientRequest,scope=Write,name=Unavailables,attr=OneMinuteRate", metric.GAUGE},
	"provider.rangeSliceTimeoutsPerSecond":      {"org.apache.cassandra.metrics:type=ClientRequest,scope=RangeSlice,name=Timeouts,attr=OneMinuteRate", metric.GAUGE},
	"provider.rangeSliceUnavalablesPerSecond":   {"org.apache.cassandra.metrics:type=ClientRequest,scope=RangeSlice,name=Unavailables,attr=OneMinuteRate", metric.GAUGE},
	"provider.commitLogCompletedTasksPerSecond": {"org.apache.cassandra.metrics:type=CommitLog,name=CompletedTasks,attr=Value", metric.COUNTER},
	"provider.commitLogPendindTasks":            {"org.apache.cassandra.metrics:type=CommitLog,name=PendingTasks,attr=Value", metric.GAUGE},
	"provider.commitLogTotalSize":               {"org.apache.cassandra.metrics:type=CommitLog,name=TotalCommitLogSize,attr=Value", metric.GAUGE},
	"software.version":                          {"version", metric.ATTRIBUTE},
}

// The patterns used to get all the beans needed for the metrics defined above
var jmxPatterns = []string{
	"org.apache.cassandra.metrics:type=ClientRequest,scope=*,name=Latency",
	"org.apache.cassandra.metrics:type=ClientRequest,scope=*,name=Timeouts",
	"org.apache.cassandra.metrics:type=ClientRequest,scope=*,name=Unavailables",
	"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=*,name=ActiveTasks",
	"org.apache.cassandra.metrics:type=ThreadPools,path=request,scope=*,name=PendingTasks",
	"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=*,name=ActiveTasks",
	"org.apache.cassandra.metrics:type=ThreadPools,path=internal,scope=*,name=PendingTasks",
	"org.apache.cassandra.metrics:type=DroppedMessage,scope=*,name=Dropped",
	"org.apache.cassandra.metrics:type=ColumnFamily,name=LiveSSTableCount",
	"org.apache.cassandra.metrics:type=Storage,name=TotalHints",
	"org.apache.cassandra.metrics:type=Storage,name=TotalHintsInProgress",
	"org.apache.cassandra.metrics:type=Cache,scope=*,name=*",
	"org.apache.cassandra.metrics:type=CommitLog,name=*",
}

func getMetrics() (map[string]interface{}, error) {
	err := jmx.Open(args.Hostname, strconv.Itoa(args.Port), args.Username, args.Password)
	if err != nil {
		return nil, err
	}
	defer jmx.Close()

	metrics := make(map[string]interface{})
	re, _ := regexp.Compile(",?attr=.*$")

	for _, query := range jmxPatterns {
		query = re.ReplaceAllString(query, "")
		q, err := jmx.Query(query)
		if err != nil {
			return nil, err
		}
		for k, v := range q {
			metrics[k] = v
		}
	}

	metrics["version"], err = getVersion()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func getVersion() (string, error) {
	cmd := exec.Command(
		"/usr/bin/nodetool",
		fmt.Sprintf("--username=%s", args.Username),
		fmt.Sprintf("--password=%s", args.Password),
		fmt.Sprintf("--host=%s", args.Hostname),
		fmt.Sprintf("--port=%d", args.Port),
		"version",
	)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("Cant fetch Cassandra version")
	}

	parts := strings.Split(string(output), ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("Cant fetch Cassandra version")
	}

	return strings.TrimSpace(parts[1]), nil
}

func populateMetrics(sample *metric.MetricSet, metrics map[string]interface{}) {
	for metricName, metricConf := range metricsDefinition {
		rawSource := metricConf[0]
		metricType := metricConf[1].(metric.SourceType)

		var rawMetric interface{}
		var ok bool

		switch source := rawSource.(type) {
		case string:
			rawMetric, ok = metrics[source]
		case func(map[string]interface{}) (float64, bool):
			rawMetric, ok = source(metrics)
		default:
			log.Warn("Invalid raw source metric for %s", metricName)
			continue
		}

		if !ok {
			log.Warn("Can't find raw metrics in results for %s", metricName)
			continue
		}

		err := sample.AddMetric(metricName, rawMetric, metricType)
		if err != nil {
			log.Warn("Error setting value: %s", err)
			continue
		}
	}
}
