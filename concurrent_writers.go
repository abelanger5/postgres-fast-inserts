package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/abelanger5/postgres-fast-inserts/internal/cmdutils"
	"github.com/abelanger5/postgres-fast-inserts/internal/dbsqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/cobra"
)

// concurrentCmd represents the concurrent command
var concurrentCmd = &cobra.Command{
	Use:   "concurrent",
	Short: "concurrent demonstrates inserts with multiple concurrent writers.",
}

var concurrentSingletonCmd = &cobra.Command{
	Use:   "singleton",
	Short: "singleton performs inserts by writing 1 row at a time.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runConcurrentSingleton(ctx)
	},
}

var concurrentBatchCmd = &cobra.Command{
	Use:   "batch",
	Short: "batch performs inserts by writing n rows within a single tx in a single database trip.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runConcurrentBatch(ctx)
	},
}

var concurrentCopyFromCmd = &cobra.Command{
	Use:   "copyfrom",
	Short: "copyfrom performs inserts by writing n rows within a single tx in a single database trip with a copy from strategy.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runConcurrentCopyfrom(ctx)
	},
}

var concurrentRowsCount int
var concurrentWritersCount int

func init() {
	rootCmd.AddCommand(concurrentCmd)

	concurrentCmd.AddCommand(concurrentSingletonCmd)
	concurrentCmd.AddCommand(concurrentBatchCmd)
	concurrentCmd.AddCommand(concurrentCopyFromCmd)

	concurrentCmd.PersistentFlags().IntVarP(
		&concurrentRowsCount,
		"count",
		"c",
		1000,
		"number of rows to insert",
	)

	concurrentCmd.PersistentFlags().IntVarP(
		&concurrentWritersCount,
		"writers",
		"w",
		10,
		"number of concurrent writers",
	)
}

func runConcurrentSingleton(ctx context.Context) {
	start := time.Now()
	reporter := NewReporter()

	wg := sync.WaitGroup{}

	for i := 0; i < concurrentWritersCount; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			batchSize := concurrentRowsCount / concurrentWritersCount
			remainder := concurrentRowsCount % concurrentWritersCount

			if i < remainder {
				batchSize++
			}

			for j := 0; j < batchSize; j++ {
				startTask := time.Now()
				payload := generateJSONPayload()

				task := dbsqlc.InsertTaskSingletonParams{
					Args: payload,
					IdempotencyKey: pgtype.Text{
						String: uuid.NewString(),
						Valid:  true,
					},
				}

				if _, err := queries.InsertTaskSingleton(ctx, pool, task); err != nil {
					log.Fatalf("could not create task: %v", err)
				}

				reporter.RecordTask(time.Since(startTask))
				reporter.RecordBatch()
			}
		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)
	reporter.Print(elapsed)
}

func runConcurrentBatch(ctx context.Context) {
	start := time.Now()

	wg := sync.WaitGroup{}
	count := 0
	countMu := sync.Mutex{}

	for i := 0; i < concurrentWritersCount; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			batchSize := concurrentRowsCount / concurrentWritersCount
			remainder := concurrentRowsCount % concurrentWritersCount

			if i < remainder {
				batchSize++
			}

			tasks := []dbsqlc.InsertTasksBatchParams{}

			for j := 0; j < batchSize; j++ {
				payload := generateJSONPayload()

				tasks = append(tasks, dbsqlc.InsertTasksBatchParams{
					Args: payload,
					IdempotencyKey: pgtype.Text{
						String: uuid.NewString(),
						Valid:  true,
					},
				})
			}

			res := queries.InsertTasksBatch(ctx, pool, tasks)

			if err := res.Close(); err != nil {
				log.Fatalf("could not create tasks batch: %v", err)
			}

			countMu.Lock()
			count += len(tasks)
			countMu.Unlock()
		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)

	log.Printf("Inserted %d rows in %s", count, elapsed)
}

func runConcurrentCopyfrom(ctx context.Context) {
	start := time.Now()

	wg := sync.WaitGroup{}
	count := 0
	countMu := sync.Mutex{}

	for i := 0; i < concurrentWritersCount; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			batchSize := concurrentRowsCount / concurrentWritersCount
			remainder := concurrentRowsCount % concurrentWritersCount

			if i < remainder {
				batchSize++
			}

			tasks := []dbsqlc.InsertTasksCopyFromParams{}

			for j := 0; j < batchSize; j++ {
				payload := generateJSONPayload()

				tasks = append(tasks, dbsqlc.InsertTasksCopyFromParams{
					Args: payload,
					IdempotencyKey: pgtype.Text{
						String: uuid.NewString(),
						Valid:  true,
					},
				})
			}

			_, err := queries.InsertTasksCopyFrom(ctx, pool, tasks)

			if err != nil {
				log.Fatalf("could not create tasks copyfrom: %v", err)
			}

			countMu.Lock()
			count += len(tasks)
			countMu.Unlock()
		}(i)
	}

	wg.Wait()

	elapsed := time.Since(start)

	log.Printf("Inserted %d rows in %s", count, elapsed)
}
