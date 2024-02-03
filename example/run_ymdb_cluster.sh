#! /bin/bash
mkdir wal && cd wal
mkdir store1 && mkdir store2 && mkdir store3 && mkdir store4 && mkdir store5  && mkdir store6  && mkdir store7  && mkdir store8  && mkdir store9
mkdir restore1 && mkdir restore2 && mkdir restore3 && mkdir restore4 && mkdir restore5  && mkdir restore6  && mkdir restore7  && mkdir restore8  && mkdir restore9

cd ..
mkdir ymdb-cluster && cd ymdb-cluster
mkdir nodeA && mkdir nodeB  && mkdir nodeC  && mkdir nodeD && mkdir nodeE && mkdir nodeF && mkdir nodeJ && mkdir nodeH && mkdir nodeI

cd ..

go build main.go
nohup ./main --raft_bootstrap --raft_id=nodeA --address=localhost:50051 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store1 --restore_file_path ./wal/restore1 &
nohup ./main --raft_id=nodeB --address=localhost:50052 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store2 --restore_file_path ./wal/restore2 &
nohup ./main --raft_id=nodeC --address=localhost:50053 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store3 --restore_file_path ./wal/restore3 &
go install github.com/Jille/raftadmin/cmd/raftadmin@latest
sleep 5s
raftadmin localhost:50051 add_voter nodeB localhost:50052 0
raftadmin --leader multi:///localhost:50051,localhost:50052 add_voter nodeC localhost:50053 0

nohup ./main --raft_bootstrap --raft_id=nodeD --address=localhost:50054 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store4 --restore_file_path ./wal/restore4 &
nohup ./main --raft_id=nodeE --address=localhost:50055 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store5 --restore_file_path ./wal/restore5 &
nohup ./main --raft_id=nodeF --address=localhost:50056 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store6 --restore_file_path ./wal/restore6 &
sleep 5s
raftadmin localhost:50054 add_voter nodeE localhost:50055 0
raftadmin --leader multi:///localhost:50054,localhost:50055 add_voter nodeF localhost:50056 0

nohup ./main --raft_bootstrap --raft_id=nodeJ --address=localhost:50057 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store7 --restore_file_path ./wal/restore7 &
nohup ./main --raft_id=nodeH --address=localhost:50058 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store8 --restore_file_path ./wal/restore8 &
nohup ./main --raft_id=nodeI --address=localhost:50059 --raft_data_dir ./ymdb-cluster --store_file_path ./wal/store9 --restore_file_path ./wal/restore9 &
sleep 5s
raftadmin localhost:50057 add_voter nodeH localhost:50058 0
raftadmin --leader multi:///localhost:50057,localhost:50058 add_voter nodeI localhost:50059 0