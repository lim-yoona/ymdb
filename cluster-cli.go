package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	log2 "log"

	_ "github.com/Jille/grpc-multi-resolver"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	cli "github.com/lim-yoona/ymdb/interact/client"
	pb "github.com/lim-yoona/ymdb/proto"
	"github.com/lim-yoona/ymdb/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
)

var Node2IP map[string]string = make(map[string]string, 9)
var Partition map[int]map[string]bool = make(map[int]map[string]bool, 3)
var Partition2Client map[int]pb.ExampleClient = make(map[int]pb.ExampleClient, 3)

func main() {
	cst := cli.InitNodeList(Node2IP)
	cli.InitPartition(Partition)
	cli.ListNode(Node2IP)
	cli.ListPartitions(Partition)
	serviceConfig := `{"healthCheckConfig": {"serviceName": "Example"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}
	for i := 0; i < len(Partition); i++ {
		target := "multi:///"
		for k, _ := range Partition[i] {
			target = target + Node2IP[k] + ","
		}
		target = target[:len(target)-1]
		conn, err := grpc.Dial(target,
			grpc.WithDefaultServiceConfig(serviceConfig), grpc.WithInsecure(),
			grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
			grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
			grpc.WithBlock())
		fmt.Println(conn.GetState().String())
		if err != nil {
			log2.Fatalf("dialing failed: %v", err)
		}
		defer conn.Close()
		c := pb.NewExampleClient(conn)

		Partition2Client[i] = c
	}

	go func() {
		for {
			var data string
			reader := bufio.NewReader(os.Stdin)
			line, _, _ := reader.ReadLine()
			data = string(line)
			// 分隔出命令和参数
			split := strings.Split(data, " ")
			log.Info().Msgf("Command will be run: %s", split)
			switch split[0] {
			case util.PUTCOMMAND:
				putMsg := &util.Put{
					Key:   split[1],
					Value: split[2],
				}
				marshal, _ := json.Marshal(putMsg)
				putMessage := &pb.SendMsgRequest{
					Id:      util.PUTID,
					DataLen: uint32(len(marshal)),
					Data:    marshal,
				}
				locate := cst.LocateKey([]byte(putMsg.Key)).String()
				var partition int
				for pn, p := range Partition {
					if p[locate] {
						partition = pn
					}
				}
				c := Partition2Client[partition]
				log.Info().Msgf("[Partition] >>>  %s stored in partition %d on node %s", putMsg.Key, partition, locate)
				_, err := c.SendMsg(context.Background(), putMessage)
				if err != nil {
					log2.Fatalf("AddWord RPC failed: %v", err)
				}
				break
			case util.GETCOMMAND:
				getMsg := &util.Other{
					Data: split[1],
				}
				marshal, _ := json.Marshal(getMsg)
				getMessage := &pb.SendMsgRequest{
					Id:      util.GETID,
					DataLen: uint32(len(marshal)),
					Data:    marshal,
				}
				locate := cst.LocateKey([]byte(getMsg.Data)).String()
				var partition int
				for pn, p := range Partition {
					if p[locate] {
						partition = pn
					}
				}
				c := Partition2Client[partition]
				log.Info().Msgf("[Partition] >>> %s stored in partition %d on node %s", getMsg.Data, partition, locate)
				res, err := c.SendMsg(context.Background(), getMessage)
				if err != nil {
					log2.Fatalf("AddWord RPC failed: %v", err)
				}
				var revResule util.Other
				json.Unmarshal(res.GetData(), &revResule)
				fmt.Println(revResule.Data)
				break
			case util.DELETECOMMAND:
				deleteMsg := &util.Other{
					Data: split[1],
				}
				marshal, _ := json.Marshal(deleteMsg)
				deleteMessage := &pb.SendMsgRequest{
					Id:      util.DELETEID,
					DataLen: uint32(len(marshal)),
					Data:    marshal,
				}
				locate := cst.LocateKey([]byte(deleteMsg.Data)).String()
				var partition int
				for pn, p := range Partition {
					if p[locate] {
						partition = pn
					}
				}
				c := Partition2Client[partition]
				log.Info().Msgf("[Partition] >>> %s stored in partition %d on node %s", deleteMsg.Data, partition, locate)
				_, err := c.SendMsg(context.Background(), deleteMessage)
				if err != nil {
					log2.Fatalf("AddWord RPC failed: %v", err)
				}
				break
			}
		}
	}()

	select {}
}
