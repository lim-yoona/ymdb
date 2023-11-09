package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/lim-yoona/tcpack"
	"github/lim-yoona/tinyKVStore/util"
	"net"
	"os"
	"strings"
)

func main() {
	address, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8099")
	if err != nil {
		fmt.Println("create address failed")
	}
	fmt.Println(*address)
	tcpConn, err := net.DialTCP("tcp4", nil, address)
	if err != nil {
		fmt.Println("create tcpconn failed", err)
	}
	mp := tcpack.NewMsgPack(8, tcpConn)
	go func() {
		for {
			var data string
			reader := bufio.NewReader(os.Stdin)
			line, _, _ := reader.ReadLine()
			data = string(line)
			// 分隔出命令和参数
			split := strings.Split(data, " ")
			fmt.Println(split)
			switch split[0] {
			case "put":
				putMsg := &util.Put{
					Key:   split[1],
					Value: split[2],
				}
				marshal, _ := json.Marshal(putMsg)
				putSendMsg := tcpack.NewMessage(1, uint32(len(marshal)), marshal)
				mp.Pack(putSendMsg)
				break
			case "get":
				getMsg := &util.Other{
					Key: split[1],
				}
				marshal, _ := json.Marshal(getMsg)
				getSendMsg := tcpack.NewMessage(2, uint32(len(marshal)), marshal)
				mp.Pack(getSendMsg)
				break
			case "delete":
				deleteMsg := &util.Other{
					Key: split[1],
				}
				marshal, _ := json.Marshal(deleteMsg)
				deleteSendMsg := tcpack.NewMessage(3, uint32(len(marshal)), marshal)
				mp.Pack(deleteSendMsg)
				break
			}
		}
	}()

	go func() {
		for {
			revMsg, _ := mp.Unpack()
			if revMsg.GetMsgId() == 1 {
				var revResule util.Other
				json.Unmarshal(revMsg.GetMsgData(), &revResule)
				fmt.Println(revResule.Key)
			}
		}
	}()
	select {}
}
