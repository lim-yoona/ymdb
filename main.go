package main

import (
	"encoding/json"
	"fmt"
	"github.com/lim-yoona/tcpack"
	"github/lim-yoona/tinyKVStore/db"
	"github/lim-yoona/tinyKVStore/util"
	"net"
)

func main() {
	yoonaDB := db.NewDB()
	address, err := net.ResolveTCPAddr("tcp4", ":8099")
	if err != nil {
		fmt.Println("create address failed")
	}
	fmt.Println(*address)
	listener, err := net.ListenTCP("tcp4", address)
	if err != nil {
		fmt.Println("create listener failed")
	}
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("return conn failed")
		}
		mp := tcpack.NewMsgPack(8, tcpConn)
		go func() {
			for {
				fmt.Println("server for")
				msg, err := mp.Unpack()
				fmt.Println(msg)
				if err == nil {
					fmt.Println("read msg failed:", err)
					switch msg.GetMsgId() {
					case uint32(1):
						var putMsg util.Put
						json.Unmarshal(msg.GetMsgData(), &putMsg)
						yoonaDB.Put(putMsg.Key, putMsg.Value)
						fmt.Println(putMsg)
						break
					case uint32(2):
						var getMsg util.Other
						json.Unmarshal(msg.GetMsgData(), &getMsg)
						get := yoonaDB.Get(getMsg.Key)
						fmt.Println(getMsg)
						sendMsg := &util.Other{
							Key: get,
						}
						marshal, _ := json.Marshal(sendMsg)
						msgSend := tcpack.NewMessage(1, uint32(len(marshal)), marshal)
						mp.Pack(msgSend)
						break
					case uint32(3):
						var deleteMsg util.Other
						json.Unmarshal(msg.GetMsgData(), &deleteMsg)
						yoonaDB.Delete(deleteMsg.Key)
						fmt.Println(deleteMsg)
						break
					}
				}
			}
		}()
	}
}
