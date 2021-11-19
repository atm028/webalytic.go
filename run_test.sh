export REDIS_HOST=10.15.80.36
export REDIS_STREAM=collector-stream
export REDIS_HANDLER=log-handlers
export REDIS_GROUP=log-handlers
export CONSUL_ENABLED=true
export CONSUL_ADDRESS=10.15.80.36:8500
export CLICKHOUSE_HOST=10.15.80.36
export CLICKHOUSE_DBNAME=webalytic
export CLICKHOUSE_FLUSH_INTERVAL=1000
export CLICKHOUSE_FLUSH_LIMIT=500
export CLICKHOUSE_HTTP_PORT=8123
export CLICKHOUSE_SERVICE_PORT=9000

go test -v -timeout 30s --run ^TestClickhouse$ ./test/common/...
