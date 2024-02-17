# ymdb
简体中文 | [English](README.md)

_**ymdb** 是一个简易的分布式键值存储系统。_ 

**ymdb** 在内存中维护一个跳表用于加快键的检索，并将值存储在磁盘上的WAL文件中。 

你可以在一个集群上运行 **ymdb** ， **分布式 ymdb** 使用 `raft` 协议来实现 **分布式共识** ，同时使用一致性哈希算法来实现 **数据分区和负载均衡**。 

**ymdb** 也保证 **崩溃一致性**。

## 架构

[![架构](https://s11.ax1x.com/2024/02/17/pFJnTYQ.jpg)](https://imgse.com/i/pFJnTYQ)

## 配置
在运行 ymdb 之前, 你需要准备一些目录:

1. wal 文件目录: 你可以在 `./config/ymDB.yaml` 中设置 wal 文件目录或者在执行 `main` 时通过 flags 设置 wal 文件目录(例如 `./main --store_file_path ./wal/store --restore_file_path ./wal/restore`)。  
2. raft 数据目录: 你可以在执行 `main` 时通过 flags 设置 raft 数据目录(例如 `./main --raft_data_dir ./ymdb-cluster`)。
   
了解目录组织和设置的更多细节, 请参考 [example](https://github.com/lim-yoona/ymdb/tree/main/example) 文件夹下的示例。

## 使用

## 集群上的 ymdb
[example](https://github.com/lim-yoona/ymdb/tree/main/example) 文件夹下提供了在集群上运行 **ymdb** 的示例:

在 ymdb 项目文件夹下执行如下脚本来启动一个 **ymdb 集群** :  
```shell
./example/run_ymdb_cluster.sh
```

现在你在一个拥有 3 个分区的 9 节点集群上运行了 ymdb 。

然后你可以执行如下脚本运行一个集群客户端来与 **ymdb 集群** 交互：
```shell
./example/run_cluster_client.sh
```

## 单机上的 ymdb
### 通过 Docker 镜像

首先，使用如下命令来拉取 **ymdb 服务端** 的 Docker 镜像：

```shell
docker pull yoonamessi/ymdb:0.1
```

然后，使用如下命令来运行一个 Docker 容器：

```shell
docker run -v ${your_host_store_path}:${path_in_ymDB.yaml} -v ${your_host_restore_path}:${path_in_ymDB.yaml} -p ${host_port}:${port_in_ymDB.yaml} -d  ymdb:0.1
```
这是一个示例:  
```shell
docker run -v /root/ymdbdata/walDir/store:/root/ymdb/walDir/store -v /root/ymdbdata/walDir/restore:/root/ymdb/walDir/restore -p 8099:8099 -d  ymdb:0.1
```
运行如下命令来启动一个 **ymdb 客户端** (注意，客户端连接的端口必须与容器向主机暴露的端口一致):  
```shell
go run ymDB-cli.go
```

现在你就可以通过 **ymdb 客户端** 来操作 ymdb 了。


## 目前支持的操作
使用 `put [key] [value]` 来存储一个 KV 对。  
使用 `get [key]` 来获取 key 对应的 value 。  
使用 `delete [key]` 来删除一个 KV 对.



## 单机上的压测
基准测试的结果或许有点奇怪，因为目前的设计中 put 不必等待 ymdb 返回，而 get 需要获得查询结果，而且通信中也有开销。
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
