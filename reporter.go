package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Reporter tracks metrics for task execution
type Reporter struct {
	taskCount    int
	numBatches   int
	totalLatency time.Duration
	mu           sync.Mutex
}

// ReportData represents the data for JSON output
type ReportData struct {
	TaskCount    int     `json:"taskCount"`
	TotalTime    string  `json:"totalTime"`
	AvgLatency   string  `json:"avgLatency"`
	Throughput   float64 `json:"throughput"`
	NumBatches   int     `json:"numBatches"`
	AvgBatchSize int     `json:"avgBatchSize"`
}

// NewReporter creates a new Reporter instance
func NewReporter() *Reporter {
	return &Reporter{
		taskCount:    0,
		totalLatency: 0,
		mu:           sync.Mutex{},
	}
}

// RecordTask records a task execution with its latency
func (r *Reporter) RecordTask(latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.taskCount++
	r.totalLatency += latency
}

// RecordBatch records a batch execution
func (r *Reporter) RecordBatch() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.numBatches++
}

// Print outputs a report of execution metrics to the console
func (r *Reporter) Print(elapsed time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var avgLatency time.Duration
	if r.taskCount > 0 {
		avgLatency = r.totalLatency / time.Duration(r.taskCount)
	}

	throughput := float64(r.taskCount) / elapsed.Seconds()

	// Calculate average batch size safely
	avgBatchSize := 0
	if r.numBatches > 0 {
		avgBatchSize = r.taskCount / r.numBatches
	}

	if jsonOutput {
		// Output as JSON
		report := ReportData{
			TaskCount:    r.taskCount,
			TotalTime:    elapsed.String(),
			AvgLatency:   avgLatency.String(),
			Throughput:   throughput,
			NumBatches:   r.numBatches,
			AvgBatchSize: avgBatchSize,
		}

		jsonBytes, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Printf("Error creating JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonBytes))
	} else {
		// Output as formatted text
		fmt.Printf("==== Execution Report ====\n")
		fmt.Printf("Total tasks executed: %d\n", r.taskCount)
		fmt.Printf("Total time: %s\n", elapsed)
		fmt.Printf("Average DB write latency: %s\n", avgLatency)
		fmt.Printf("Throughput: %.2f rows/second\n", throughput)
		fmt.Printf("Number of batches: %d\n", r.numBatches)
		fmt.Printf("Average batch size: %d\n", avgBatchSize)
		fmt.Printf("========================\n")
	}
}
