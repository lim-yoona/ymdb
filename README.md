# ymdb
ymdb is a simple KV storage system that supports storing KV pairs of string types. It maintains a skiplist in memory to speed up key retrieval and stores values in WAL files on disk. ymdb also supports crash consistency.  

## Config
Before using, you need to modify some config items in `./config/ymDB.yaml`.

## Usage

### By using source code
Run the following command to start a ymdb database server:  
```shell
go run main.go
```
Run the following command to start a ymdb database client:  
```shell
go run ymDB-cli.go
```
Then you can manipulate ymdb through the database client.  

### By using Docker image
Firstly, use the following command to pull the Docker image of the ymdb server:  
```shell
docker pull yoonamessi/ymdb:0.1
```
Then, use the following command to run a Docker container:  
```shell
docker run -v ${your_host_store_path}:${path_in_ymDB.yaml} -v ${your_host_restore_path}:${path_in_ymDB.yaml} -p ${host_port}:${port_in_ymDB.yaml} -d  ymdb:0.1
```
This is an example:  
```shell
docker run -v /root/ymdbdata/walDir/store:/root/ymdb/walDir/store -v /root/ymdbdata/walDir/restore:/root/ymdb/walDir/restore -p 8099:8099 -d  ymdb:0.1
```
Run the following command to start a ymdb database client(Note that the port connected by the client must match the port exposed by the container to the host):  
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