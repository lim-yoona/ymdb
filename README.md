# tinyKVStore
ymdb is a simple KV storage system that supports storing KV pairs of string types. It maintains a skiplist in memory to speed up key retrieval and stores values in WAL files on disk. ymdb also supports crash consistency.  

## Config
Before using, you need to modify some config items in `./config/ymDB.yaml`.

## Usage
Run the following command to start a ymdb database server:  
```shell
go run main.go
```
Run the following command to start a ymdb database client:  
```shell
go run ymDB-cli.go
```
Then you can manipulate ymdb through the database client.  

## Currently supported commands
Using `put [key] [value]` to store a KV pair to the database.  
Using `get [key]` to get the value of the key.  
Using `delete [key]` to delete a KV pair.

## TODO
ymdb is under development and future plans are as follows:  
- Support high concurrency.
- Support distributed KV storage.

## Benchmark
The benchmark result is a bit weird because put doesn't have to wait for ymdb to return, get needs to get the query result, and there's also the overhead involved in the communication.  
```shell
goos: windows
goarch: amd64
pkg: github.com/lim-yoona/ymdb/benchmark
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkPutGet
BenchmarkPutGet/put
BenchmarkPutGet/put-16            339727              2994 ns/op            2392 B/op         40 allocs/op
BenchmarkPutGet/get
BenchmarkPutGet/get-16            915457              2891 ns/op            2186 B/op         37 allocs/op
PASS
```