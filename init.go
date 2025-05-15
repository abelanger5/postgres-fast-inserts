package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/abelanger5/postgres-fast-inserts/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool
var queries *dbsqlc.Queries
var maxConns int
var continuousWritersCount int
var benchmarkDuration time.Duration
var flushInterval time.Duration
var batchSize int
var channelBufferSize int
var jsonOutput bool
var withAssociatedData bool

func init() {
	rootCmd.PersistentFlags().IntVarP(&maxConns, "max-conns", "m", 20, "maximum number of connections to the database")

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	config, err := pgxpool.ParseConfig(dbUrl)

	if err != nil {
		log.Fatal("could not parse DATABASE_URL: %w", err)
	}

	config.MaxConns = int32(maxConns)

	pool, err = pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatal("could not create connection pool: %w", err)
	}

	queries = dbsqlc.New()

	rootCmd.AddCommand(continuousCmd)

	rootCmd.PersistentFlags().BoolVarP(
		&jsonOutput,
		"json",
		"j",
		false,
		"output in JSON format",
	)

	continuousCmd.AddCommand(continuousSingletonCmd)
	continuousCmd.AddCommand(continuousBatchCmd)
	continuousCmd.AddCommand(continuousCopyFromCmd)
	continuousCmd.AddCommand(continuousPingCmd)

	continuousCmd.PersistentFlags().IntVarP(
		&continuousWritersCount,
		"writers",
		"w",
		10,
		"number of continuous writers",
	)

	if continuousWritersCount > maxConns {
		log.Fatalf("number of writers (%d) cannot be greater than max connections (%d). increase max connections via the --max-conns flag", continuousWritersCount, maxConns)
	}

	continuousCmd.PersistentFlags().DurationVarP(
		&benchmarkDuration,
		"duration",
		"d",
		10*time.Second,
		"duration to run the benchmark (e.g. 10s, 1m, 2h)",
	)

	continuousCmd.PersistentFlags().IntVarP(
		&batchSize,
		"batch-size",
		"b",
		100,
		"batch size for batch and copyfrom strategies",
	)

	continuousCmd.PersistentFlags().IntVarP(
		&channelBufferSize,
		"buffer",
		"",
		batchSize*continuousWritersCount,
		"size of the buffer channel for tasks",
	)

	continuousCmd.PersistentFlags().BoolVar(
		&withAssociatedData,
		"with-associated-data",
		false,
		"insert associated data with the task",
	)

	continuousCmd.PersistentFlags().DurationVarP(
		&flushInterval,
		"flush-interval",
		"f",
		10*time.Millisecond,
		"interval to flush the buffer",
	)
}
