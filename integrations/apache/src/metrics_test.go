package main

import (
	"bufio"
	"strings"
	"testing"
)

var testApacheStatus = `Total Accesses: 66
Total kBytes: 73
Uptime: 31006
ReqPerSec: .00212862
BytesPerSec: 2.41089
BytesPerReq: 1132.61
BusyWorkers: 1
IdleWorkers: 4
Scoreboard: _W___......_CDCDII.II......KKKKKGG................__R_W.....S.....LS
`

var testApacheStatusWrongLinesFormat = `
Random text
Random text

`

var testApacheStatusEmpty = ``

func TestGetRawMetrics(t *testing.T) {
	rawMetrics, err := getRawMetrics(bufio.NewReader(strings.NewReader(testApacheStatus)))

	if len(rawMetrics) != 9 {
		t.Error()
	}
	if rawMetrics["Total Accesses"] != 66 {
		t.Error()
	}
	if rawMetrics["Uptime"] != 31006 {
		t.Error()
	}
	if rawMetrics["ReqPerSec"] != 0.00212862 {
		t.Error()
	}
	if rawMetrics["BytesPerSec"] != 2.41089 {
		t.Error()
	}
	if rawMetrics["IdleWorkers"] != 4 {
		t.Error()
	}
	if rawMetrics["BusyWorkers"] != 1 {
		t.Error()
	}
	if rawMetrics["Total kBytes"] != 73 {
		t.Error()
	}
	if rawMetrics["BytesPerReq"] != 1132.61 {
		t.Error()
	}
	if rawMetrics["Scoreboard"] != "_W___......_CDCDII.II......KKKKKGG................__R_W.....S.....LS" {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestGetMetricsInvalidData(t *testing.T) {
	rawMetrics, err := getRawMetrics(bufio.NewReader(strings.NewReader(testApacheStatusWrongLinesFormat)))

	if err == nil {
		t.Error()
	}
	if rawMetrics != nil {
		t.Error()
	}
}

func TestGetMetricsEmptyData(t *testing.T) {
	rawMetrics, err := getRawMetrics(bufio.NewReader(strings.NewReader(testApacheStatusEmpty)))

	if err == nil {
		t.Error()
	}
	if rawMetrics != nil {
		t.Error()
	}
}

func TestGetWorkerStatus(t *testing.T) {
	metrics := map[string]interface{}{
		"Scoreboard": "_W___......_DDII.II......KKKKKGG................__R_W.....S.....LS",
	}
	writingWorkersNumber, ok := getWorkerStatus("W")(metrics)
	if ok != true {
		t.Error()
	}
	if writingWorkersNumber != float64(2) {
		t.Error()
	}

	closingWorkersNumber, ok := getWorkerStatus("C")(metrics)
	if ok != true {
		t.Error()
	}
	if closingWorkersNumber != float64(0) {
		t.Error()
	}
}

func TestGetWorkerStatusInvalidDataKey(t *testing.T) {
	metrics := map[string]interface{}{
		"Total kBytes": "_W___......_CDCDII.II......KKKKKGG................__R_W.....S.....LS",
	}
	closingWorkersNumber, ok := getWorkerStatus("C")(metrics)
	if ok != false {
		t.Error()
	}
	if closingWorkersNumber != float64(0) {
		t.Error()
	}
}

func TestGetWorkerStatusInvalidDataType(t *testing.T) {
	metrics := map[string]interface{}{
		"Scoreboard": 5,
	}
	closingWorkersNumber, ok := getWorkerStatus("C")(metrics)
	if ok != false {
		t.Error()
	}
	if closingWorkersNumber != float64(0) {
		t.Error()
	}
}

func TestGetTotalWorkers(t *testing.T) {
	metrics := map[string]interface{}{
		"Scoreboard": "_W___......_DDII.II......KKKKKGG................__R_W.....S.....LS",
	}
	totalWorkersNumber, ok := getTotalWorkers(metrics)
	if ok != true {
		t.Error()
	}
	if totalWorkersNumber != float64(66) {
		t.Error()
	}
}

func TestGetTotalWorkersInvalidDataKey(t *testing.T) {
	metrics := map[string]interface{}{
		"Total kBytes": "_W___......_CDCDII.II......KKKKKGG................__R_W.....S.....LS",
	}
	totalWorkersNumber, ok := getTotalWorkers(metrics)
	if ok != false {
		t.Error()
	}
	if totalWorkersNumber != float64(0) {
		t.Error()
	}
}

func TestGetTotalWorkersInvalidDataType(t *testing.T) {
	metrics := map[string]interface{}{
		"Scoreboard": 5,
	}
	totalWorkersNumber, ok := getTotalWorkers(metrics)
	if ok != false {
		t.Error()
	}
	if totalWorkersNumber != float64(0) {
		t.Error()
	}
}

func TestGetBytes_IntData(t *testing.T) {
	metrics := map[string]interface{}{
		"Total kBytes": 67,
	}
	totalBytes, ok := getBytes(metrics)
	if ok != true {
		t.Error()
	}
	if totalBytes != float64(68608) {
		t.Error()
	}
}

func TestGetBytes_InvalidDataType(t *testing.T) {
	metrics := map[string]interface{}{
		"Total kBytes": 67.4,
	}
	totalBytes, ok := getBytes(metrics)
	if ok != false {
		t.Error()
	}
	if totalBytes != float64(0) {
		t.Error()
	}
}

func TestGetBytes_InvalidDataKey(t *testing.T) {
	metrics := map[string]interface{}{
		"TotalkBytes": 67,
	}
	totalBytes, ok := getBytes(metrics)
	if ok != false {
		t.Error()
	}
	if totalBytes != float64(0) {
		t.Error()
	}
}
