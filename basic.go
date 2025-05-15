package main

import (
	"context"
	"log"
	"time"

	"github.com/abelanger5/postgres-fast-inserts/internal/cmdutils"
	"github.com/abelanger5/postgres-fast-inserts/internal/dbsqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/cobra"
)

// basicCmd represents the seed command
var basicCmd = &cobra.Command{
	Use:   "basic",
	Short: "basic demonstrates basic strategies for fast inserts.",
}

var basicSingletonCmd = &cobra.Command{
	Use:   "singleton",
	Short: "singleton performs inserts by writing 1 row at a time.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runSingleton(ctx)
	},
}

var basicBatchCmd = &cobra.Command{
	Use:   "batch",
	Short: "batch performs inserts by writing n rows within a single tx in a single database trip.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runBatch(ctx)
	},
}

var basicBulkCmd = &cobra.Command{
	Use:   "unnest",
	Short: "unnest performs inserts by writing n rows within a single tx in a single database trip with an unnest strategy.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runUnnest(ctx)
	},
}

var basicCopyFromCmd = &cobra.Command{
	Use:   "copyfrom",
	Short: "copyfrom performs inserts by writing n rows within a single tx in a single database trip with a copy from strategy.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runCopyFrom(ctx)
	},
}

var basicCount int

func init() {
	rootCmd.AddCommand(basicCmd)

	basicCmd.AddCommand(basicSingletonCmd)
	basicCmd.AddCommand(basicBatchCmd)
	basicCmd.AddCommand(basicBulkCmd)
	basicCmd.AddCommand(basicCopyFromCmd)

	basicCmd.PersistentFlags().IntVarP(
		&basicCount,
		"count",
		"c",
		1000,
		"number of rows to insert",
	)
}

func runSingleton(ctx context.Context) {
	start := time.Now()
	reporter := NewReporter()

	for i := 0; i < basicCount; i++ {
		insertStart := time.Now()
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

		reporter.RecordBatch()
		reporter.RecordTask(time.Since(insertStart))
	}

	reporter.Print(time.Since(start))
}

func runBatch(ctx context.Context) {
	start := time.Now()
	reporter := NewReporter()

	tasks := []dbsqlc.InsertTasksBatchParams{}

	for i := 0; i < basicCount; i++ {
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

	for i := 0; i < basicCount; i++ {
		reporter.RecordTask(time.Since(start))
	}

	reporter.RecordBatch()

	reporter.Print(time.Since(start))
}

func runUnnest(ctx context.Context) {
	start := time.Now()
	reporter := NewReporter()

	taskArgs := [][]byte{}
	taskKeys := []string{}

	for i := 0; i < basicCount; i++ {
		payload := generateJSONPayload()
		taskArgs = append(taskArgs, payload)
		taskKeys = append(taskKeys, uuid.NewString())
	}

	_, err := queries.InsertTasksWithUnnest(ctx, pool, dbsqlc.InsertTasksWithUnnestParams{
		Args: taskArgs,
		Keys: taskKeys,
	})

	if err != nil {
		log.Fatalf("could not create tasks batch: %v", err)
	}

	for i := 0; i < basicCount; i++ {
		reporter.RecordTask(time.Since(start))
	}

	reporter.RecordBatch()
	reporter.Print(time.Since(start))
}

func runCopyFrom(ctx context.Context) {
	start := time.Now()
	reporter := NewReporter()

	tasks := []dbsqlc.InsertTasksCopyFromParams{}

	for i := 0; i < basicCount; i++ {
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
		log.Fatalf("could not create tasks batch: %v", err)
	}

	for i := 0; i < basicCount; i++ {
		reporter.RecordTask(time.Since(start))
	}

	reporter.RecordBatch()

	reporter.Print(time.Since(start))
}
