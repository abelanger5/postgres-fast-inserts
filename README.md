This repository accompanies [this blog post](https://docs.hatchet.run/blog/fastest-postgres-inserts).

## Benchmarks

Benchmarks can be run via `sh run-benchmarks.sh`. The following benchmarks were run on a 2023 MacBook Pro, Apple M3 Max chip, 36 GB memory, macOS 14.3.

<details><summary><strong>Results</strong></summary>

```
Running pg-inserts basic singleton: 1 connection, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 43.782110417s
Average DB write latency: 437.71µs
Throughput: 2284.04 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts basic batch: 1 connection, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 2.669160083s
Average DB write latency: 2.668188328s
Throughput: 37464.97 rows/second
Number of batches: 1
Average batch size: 100000
========================

Running pg-inserts basic copyfrom: 1 connection, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 1.580160667s
Average DB write latency: 1.577824529s
Throughput: 63284.70 rows/second
Number of batches: 1
Average batch size: 100000
========================

Running pg-inserts concurrent singleton: 10 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 9.087605458s
Average DB write latency: 907.187µs
Throughput: 11004.00 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 20 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 6.004460458s
Average DB write latency: 1.196766ms
Throughput: 16654.29 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 30 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 6.058386667s
Average DB write latency: 1.809936ms
Throughput: 16506.04 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 40 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 6.048169041s
Average DB write latency: 2.415203ms
Throughput: 16533.93 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 50 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 5.953354834s
Average DB write latency: 2.971404ms
Throughput: 16797.25 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 60 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 5.894446125s
Average DB write latency: 3.53187ms
Throughput: 16965.12 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 70 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 5.880129625s
Average DB write latency: 4.110432ms
Throughput: 17006.43 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 80 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 6.113362625s
Average DB write latency: 4.884029ms
Throughput: 16357.61 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 90 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 5.925777709s
Average DB write latency: 5.326193ms
Throughput: 16875.42 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts concurrent singleton: 100 connections, 100000 rows, no batching
==== Execution Report ====
Total tasks executed: 100000
Total time: 5.881917083s
Average DB write latency: 5.872899ms
Throughput: 17001.26 rows/second
Number of batches: 100000
Average batch size: 1
========================

Running pg-inserts continuous ping (this doesn't write any data): 30 seconds, 20 connections, batch size 100
==== Execution Report ====
Total tasks executed: 3261115
Total time: 30.010939041s
Average DB write latency: 4.721006ms
Throughput: 108664.21 rows/second
Number of batches: 49790
Average batch size: 65
========================

Running pg-inserts continuous batch: 30 seconds, 20 connections, batch size 100
==== Execution Report ====
Total tasks executed: 2399523
Total time: 30.021498625s
Average DB write latency: 42.973188ms
Throughput: 79926.82 rows/second
Number of batches: 24006
Average batch size: 99
========================

Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 100
==== Execution Report ====
Total tasks executed: 2826355
Total time: 30.676920125s
Average DB write latency: 17.986946ms
Throughput: 92132.95 rows/second
Number of batches: 39423
Average batch size: 71
========================

Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 5, flush interval 0.5ms
==== Execution Report ====
Total tasks executed: 1432368
Total time: 30.003068083s
Average DB write latency: 3.973591ms
Throughput: 47740.72 rows/second
Number of batches: 286475
Average batch size: 4
========================

Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 10, flush interval 1ms
==== Execution Report ====
Total tasks executed: 1817467
Total time: 30.0029295s
Average DB write latency: 6.157382ms
Throughput: 60576.32 rows/second
Number of batches: 181748
Average batch size: 9
========================

Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 25, flush interval 2.5ms
==== Execution Report ====
Total tasks executed: 2663036
Total time: 30.008680125s
Average DB write latency: 9.642757ms
Throughput: 88742.19 rows/second
Number of batches: 106551
Average batch size: 24
========================

Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 50, flush interval 5ms
==== Execution Report ====
Total tasks executed: 3222520
Total time: 30.014999417s
Average DB write latency: 8.248439ms
Throughput: 107363.65 rows/second
Number of batches: 79585
Average batch size: 40
========================

Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 100, flush interval 10ms
==== Execution Report ====
Total tasks executed: 2949790
Total time: 30.014832125s
Average DB write latency: 12.795711ms
Throughput: 98277.74 rows/second
Number of batches: 43068
Average batch size: 68
========================

Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 200, flush interval 20ms
==== Execution Report ====
Total tasks executed: 3350520
Total time: 30.022126292s
Average DB write latency: 14.98438ms
Throughput: 111601.69 rows/second
Number of batches: 26107
Average batch size: 128
========================
```

</details>

## Setup

**Prerequisites:**

- [Taskfile](https://taskfile.dev/)
- Go 1.24+
- Docker Compose

Run `task setup` to get everything running. This will spin up a Postgres database on port 5432, generate the relevant `sqlc`, and write the schema to the database. You might need to run it multiple times if Postgres doesn't start quickly.

Next, set the `DATABASE_URL` environment variable for all commands below:

```
export DATABASE_URL=postgresql://hatchet:hatchet@127.0.0.1:5432/hatchet
```

Build the Go binary:

```sh
go build -o ./bin/pg-inserts .
mv ./bin/pg-inserts /usr/local/bin
```

Signature:

```
$ pg-inserts -h
inserts demonstrates simple commands for testing different insert methods in Postgres.

Usage:
  inserts [flags]
  inserts [command]

Available Commands:
  basic       basic demonstrates basic strategies for fast inserts.
  completion  Generate the autocompletion script for the specified shell
  concurrent  concurrent demonstrates inserts with multiple concurrent writers.
  continuous  continuous demonstrates inserts with multiple continuous writers.
  help        Help about any command

Flags:
  -h, --help                   help for inserts
  -j, --json                   output in JSON format
  -m, --max-conns int          maximum number of connections to the database (default 20)
      --max-payload-size int   maximum size of the payload in kilobytes (default 1000)

Use "inserts [command] --help" for more information about a command.
```
