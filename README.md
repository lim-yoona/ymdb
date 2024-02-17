# ymdb
English | [简体中文](README-CN.md)

_**ymdb** is a simple distributed key-value storage system._ 

**ymdb** maintains a skip-list in memory to speed up key retrieval and stores values in WAL files on disk. 

You can run **ymdb** on a cluster, **distributed ymdb** uses the `raft` protocol for **distributed consensus** and a consistent hashing algorithm for **data partitioning and load balancing**. 

**ymdb** also supports **crash consistency**.

## Architecture

[![Architecture](https://s11.ax1x.com/2024/02/17/pFJnTYQ.jpg)](https://imgse.com/i/pFJnTYQ)

## Config
Before running ymdb, you need to prepare some directories:

1. wal file directory: you can set wal file directory in `./config/ymDB.yaml` or set wal file directory by flags when execute `main`(eg. `./main --store_file_path ./wal/store --restore_file_path ./wal/restore`).  
2. raft data directory: you can set raft data directory by flags when execute `main`(eg. `./main --raft_data_dir ./ymdb-cluster`)
   
For a more detailed directory organization and setup, see the example under [example](https://github.com/lim-yoona/ymdb/tree/main/example) folder.

## Usage

## ymdb on a cluster
An example of **ymdb** running on a cluster is provided under the [example](https://github.com/lim-yoona/ymdb/tree/main/example) folder:

Start an **ymdb cluster** by executing the following script under the ymdb project folder:  
```shell
./example/run_ymdb_cluster.sh
```
Now you have ymdb running on a nine-node cluster with three partitions.

You can then run a cluster client and interact with the **ymdb cluster** by executing the following script:  
```shell
./example/run_cluster_client.sh
```

## ymdb on a single machine
### By using Docker image
Firstly, use the following command to pull the Docker image of the **ymdb server**:  
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
Run the following command to start a **ymdb client**(Note that the port connected by the client must match the port exposed by the container to the host):  
```shell
go run ymDB-cli.go
```
Then you can manipulate ymdb through the database client.


## Currently supported commands
Using `put [key] [value]` to store a KV pair to the database.  
Using `get [key]` to get the value of the key.  
Using `delete [key]` to delete a KV pair.



## Benchmark on a single machine
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
