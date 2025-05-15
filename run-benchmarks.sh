
echo "Running benchmarks..."

task setup

echo "Running pg-inserts basic singleton: 1 connection, 100000 rows, no batching"
pg-inserts basic singleton -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts basic batch: 1 connection, 100000 rows, no batching"
pg-inserts basic batch -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts basic copyfrom: 1 connection, 100000 rows, no batching"
pg-inserts basic copyfrom -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 10 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 10 --writers 10 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 20 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 20 --writers 20 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 30 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 30 --writers 30 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 40 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 40 --writers 40 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 50 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 50 --writers 50 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 60 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 60 --writers 60 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 70 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 70 --writers 70 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 80 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 80 --writers 80 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 90 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 90 --writers 90 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts concurrent singleton: 100 connections, 100000 rows, no batching"
pg-inserts concurrent singleton --max-conns 100 --writers 100 -c 100000

task reset-db &> /dev/null

echo "Running pg-inserts continuous ping (this doesn't write any data): 30 seconds, 20 connections, batch size 100"
pg-inserts continuous ping --duration 30s --batch-size 100 --max-conns 20 --writers 20

task reset-db &> /dev/null

echo "Running pg-inserts continuous batch: 30 seconds, 20 connections, batch size 100"
pg-inserts continuous batch --duration 30s --batch-size 100 --max-conns 20 --writers 20

task reset-db &> /dev/null

echo "Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 100"
pg-inserts continuous copyfrom --duration 30s --batch-size 100 --max-conns 20 --writers 20

task reset-db &> /dev/null

echo "Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 5, flush interval 0.5ms"
pg-inserts continuous copyfrom --duration 30s --batch-size 5 --max-conns 20 --writers 20 --flush-interval 500µs

task reset-db &> /dev/null

echo "Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 10, flush interval 1ms"
pg-inserts continuous copyfrom --duration 30s --batch-size 10 --max-conns 20 --writers 20 --flush-interval 1ms

task reset-db &> /dev/null

echo "Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 25, flush interval 2.5ms"
pg-inserts continuous copyfrom --duration 30s --batch-size 25 --max-conns 20 --writers 20 --flush-interval 2500µs

task reset-db &> /dev/null

echo "Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 50, flush interval 5ms"
pg-inserts continuous copyfrom --duration 30s --batch-size 50 --max-conns 20 --writers 20 --flush-interval 5ms

task reset-db &> /dev/null

echo "Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 100, flush interval 10ms"
pg-inserts continuous copyfrom --duration 30s --batch-size 100 --max-conns 20 --writers 20 --flush-interval 10ms

task reset-db &> /dev/null

echo "Running pg-inserts continuous copyfrom: 30 seconds, 20 connections, batch size 200, flush interval 20ms"
pg-inserts continuous copyfrom --duration 30s --batch-size 200 --max-conns 20 --writers 20 --flush-interval 20ms

# echo "Rows in database:"
# psql 'postgresql://hatchet:hatchet@127.0.0.1:5432/hatchet' -c "SELECT COUNT(*) FROM tasks"