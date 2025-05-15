package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/abelanger5/postgres-fast-inserts/internal/cmdutils"
	"github.com/abelanger5/postgres-fast-inserts/internal/dbsqlc"
	"github.com/spf13/cobra"
)

// continuousCmd represents the continuous command
var continuousCmd = &cobra.Command{
	Use:   "continuous",
	Short: "continuous demonstrates inserts with multiple continuous writers.",
}

var continuousSingletonCmd = &cobra.Command{
	Use:   "singleton",
	Short: "singleton performs inserts by writing 1 row at a time.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runContinuousSingleton(ctx)
	},
}

var continuousBatchCmd = &cobra.Command{
	Use:   "batch",
	Short: "batch performs inserts by writing n rows within a single tx in a single database trip.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runContinuousBatch(ctx)
	},
}

var continuousCopyFromCmd = &cobra.Command{
	Use:   "copyfrom",
	Short: "copyfrom performs inserts by writing n rows within a single tx in a single database trip with a copy from strategy.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runContinuousCopyfrom(ctx)
	},
}

var continuousPingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping performs a ping to the database instead of a write, to determine baseline performance.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := cmdutils.NewInterruptContext()
		defer cancel()

		runContinuousPing(ctx)
	},
}

func runContinuousSingleton(ctx context.Context) {
	// Create a data generator
	generator := NewDataGenerator(channelBufferSize)
	reporter := NewReporter()
	semaphore := make(chan struct{}, continuousWritersCount)

	// Set up context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, benchmarkDuration)
	defer cancel()

	generator.Start(ctx)

	wg := sync.WaitGroup{}
	count := 0
	countMu := sync.Mutex{}

	start := time.Now()

	var insertFunc func(context.Context, dbsqlc.InsertTaskSingletonParams) error

	insertFunc = insertSingletonBasic

	if withAssociatedData {
		insertFunc = insertSingletonBasicWithAssociatedData
	}

outer:
	for {
		select {
		case <-timeoutCtx.Done():
			break outer
		case task, ok := <-generator.Tasks():
			if !ok {
				break outer
			}

			semaphore <- struct{}{}

			wg.Add(1)

			go func(task TaskParams) {
				defer wg.Done()
				defer func() { <-semaphore }()

				startTime := time.Now()

				singletonCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err := insertFunc(singletonCtx, dbsqlc.InsertTaskSingletonParams{
					Args:           task.Args,
					IdempotencyKey: task.IdempotencyKey,
				})

				// Record latency for this task
				latency := time.Since(startTime)
				reporter.RecordTask(latency)
				reporter.RecordBatch()

				if err != nil {
					log.Printf("could not create task: %v", err)
					return
				}
			}(task)

			countMu.Lock()
			count++
			countMu.Unlock()
		}
	}

	// Wait for all workers to finish
	wg.Wait()

	reporter.Print(time.Since(start))
}

func runContinuousBatch(ctx context.Context) {
	// Create a reporter
	reporter := NewReporter()

	insertFunc := insertBatch

	if withAssociatedData {
		insertFunc = insertBatchWithAssociatedData
	}

	writeFunc := func(tasks []dbsqlc.InsertTasksBatchParams) ([]*dbsqlc.Task, error) {
		reporter.RecordBatch()

		return insertFunc(ctx, tasks)
	}

	// Create a data generator
	generator := NewDataGenerator(channelBufferSize)
	buffer := NewBuffer(ctx, writeFunc)

	// Set up context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, benchmarkDuration)
	defer cancel()

	generator.Start(ctx)

	var wg sync.WaitGroup

	start := time.Now()

outer:
	for {
		select {
		case <-timeoutCtx.Done():
			break outer
		case task, ok := <-generator.Tasks():
			if !ok {
				break outer
			}

			startTime := time.Now()

			taskWithCh, err := buffer.WriteNoWait(dbsqlc.InsertTasksBatchParams{
				Args:           task.Args,
				IdempotencyKey: task.IdempotencyKey,
			})

			if err != nil {
				log.Printf("could not buffer task: %v", err)
				return
			}

			wg.Add(1)

			go func(task TaskParams) {
				defer wg.Done()

				_, err := taskWithCh.GetResult()

				// Record latency for this task
				latency := time.Since(startTime)
				reporter.RecordTask(latency)

				if err != nil {
					log.Printf("could not create task: %v", err)
					return
				}
			}(task)
		}
	}

	// Wait for all workers to finish
	wg.Wait()

	elapsed := time.Since(start)

	// Print the report
	reporter.Print(elapsed)
}

func runContinuousPing(ctx context.Context) {
	// Create a reporter
	reporter := NewReporter()

	writeFunc := func(tasks []dbsqlc.InsertTasksBatchParams) ([]*dbsqlc.InsertTasksBatchParams, error) {
		reporter.RecordBatch()

		err := pool.Ping(ctx)

		if err != nil {
			log.Fatalf("could not ping database: %v", err)
		}

		resTasks := make([]*dbsqlc.InsertTasksBatchParams, 0, len(tasks))

		for i := 0; i < len(tasks); i++ {
			resTasks = append(resTasks, &tasks[i])
		}

		return resTasks, nil
	}

	// Create a data generator
	generator := NewDataGenerator(channelBufferSize)
	buffer := NewBuffer(ctx, writeFunc)

	// Set up context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, benchmarkDuration)
	defer cancel()

	generator.Start(ctx)

	var wg sync.WaitGroup

	start := time.Now()

outer:
	for {
		select {
		case <-timeoutCtx.Done():
			break outer
		case task, ok := <-generator.Tasks():
			if !ok {
				break outer
			}

			startTime := time.Now()

			taskWithCh, err := buffer.WriteNoWait(dbsqlc.InsertTasksBatchParams{
				Args:           task.Args,
				IdempotencyKey: task.IdempotencyKey,
			})

			if err != nil {
				log.Printf("could not buffer task: %v", err)
				return
			}

			wg.Add(1)

			go func(task TaskParams) {
				defer wg.Done()

				_, err := taskWithCh.GetResult()

				// Record latency for this task
				latency := time.Since(startTime)
				reporter.RecordTask(latency)

				if err != nil {
					log.Printf("could not create task: %v", err)
					return
				}
			}(task)
		}
	}

	// Wait for all workers to finish
	wg.Wait()

	elapsed := time.Since(start)

	// Print the report
	reporter.Print(elapsed)
}

func runContinuousCopyfrom(ctx context.Context) {
	// Create a reporter
	reporter := NewReporter()

	// note: since copyfrom doesn't return the created tasks, we just return the input
	writeFunc := func(tasks []dbsqlc.InsertTasksCopyFromParams) ([]*dbsqlc.InsertTasksCopyFromParams, error) {
		reporter.RecordBatch()

		n, err := queries.InsertTasksCopyFrom(ctx, pool, tasks)

		if err != nil {
			log.Fatalf("could not create tasks copyfrom: %v", err)
		}

		if int(n) != len(tasks) {
			log.Fatalf("could not create tasks copyfrom: expected %d, got %d", len(tasks), n)
		}

		resTasks := make([]*dbsqlc.InsertTasksCopyFromParams, 0, len(tasks))

		for i := 0; i < len(tasks); i++ {
			resTasks = append(resTasks, &tasks[i])
		}

		return resTasks, nil
	}

	// Create a data generator
	generator := NewDataGenerator(channelBufferSize)
	buffer := NewBuffer(ctx, writeFunc)

	// Set up context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, benchmarkDuration)
	defer cancel()

	generator.Start(ctx)

	var wg sync.WaitGroup

	start := time.Now()

outer:
	for {
		select {
		case <-timeoutCtx.Done():
			break outer
		case task, ok := <-generator.Tasks():
			if !ok {
				break outer
			}

			startTime := time.Now()

			taskWithCh, err := buffer.WriteNoWait(dbsqlc.InsertTasksCopyFromParams{
				Args:           task.Args,
				IdempotencyKey: task.IdempotencyKey,
			})

			if err != nil {
				log.Printf("could not buffer task: %v", err)
				return
			}

			wg.Add(1)

			go func(task TaskParams) {
				defer wg.Done()

				_, err := taskWithCh.GetResult()

				// Record latency for this task
				latency := time.Since(startTime)
				reporter.RecordTask(latency)

				if err != nil {
					log.Printf("could not create task: %v", err)
					return
				}
			}(task)
		}
	}

	// Wait for all workers to finish
	wg.Wait()

	elapsed := time.Since(start)

	// Print the report
	reporter.Print(elapsed)
}

func insertSingletonBasic(ctx context.Context, params dbsqlc.InsertTaskSingletonParams) error {
	_, err := queries.InsertTaskSingleton(ctx, pool, params)

	return err
}

func insertSingletonBasicWithAssociatedData(ctx context.Context, params dbsqlc.InsertTaskSingletonParams) error {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	task, err := queries.InsertTaskSingleton(ctx, tx, params)

	if err != nil {
		return err
	}

	err = queries.InsertTaskAssociatedData(ctx, tx, dbsqlc.InsertTaskAssociatedDataParams{
		TaskID:   task.ID,
		ArgsJson: params.Args,
	})

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return err
}

func insertBatch(ctx context.Context, tasks []dbsqlc.InsertTasksBatchParams) ([]*dbsqlc.Task, error) {
	res := queries.InsertTasksBatch(ctx, pool, tasks)

	resTasks := make([]*dbsqlc.Task, 0, len(tasks))

	res.QueryRow(func(i int, t *dbsqlc.Task, err error) {
		if err != nil {
			log.Printf("could not create task: %v", err)
			return
		}
		resTasks = append(resTasks, t)
	})

	if err := res.Close(); err != nil {
		log.Fatalf("could not create tasks batch: %v", err)
	}

	return resTasks, nil
}

func insertBatchWithAssociatedData(ctx context.Context, tasks []dbsqlc.InsertTasksBatchParams) ([]*dbsqlc.Task, error) {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	res := queries.InsertTasksBatch(ctx, tx, tasks)

	resTasks := make([]*dbsqlc.Task, 0, len(tasks))

	res.QueryRow(func(i int, t *dbsqlc.Task, err error) {
		if err != nil {
			log.Printf("could not create task: %v", err)
			return
		}

		resTasks = append(resTasks, t)
	})

	if err := res.Close(); err != nil {
		log.Fatalf("could not create tasks batch: %v", err)
	}

	args2 := make([]dbsqlc.InsertTaskAssociatedDatasBatchParams, 0, len(tasks))

	for _, task := range resTasks {
		args2 = append(args2, dbsqlc.InsertTaskAssociatedDatasBatchParams{
			TaskID:   task.ID,
			ArgsJson: task.Args,
		})
	}

	res2 := queries.InsertTaskAssociatedDatasBatch(ctx, tx, args2)

	if err := res2.Close(); err != nil {
		log.Fatalf("could not create task associated data batch: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return resTasks, err
}
