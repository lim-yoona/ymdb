package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/lim-yoona/tcpack"
	"github.com/lim-yoona/ymdb/config"
	"github.com/lim-yoona/ymdb/util"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"strings"
)

func ClientStart() {
	address, err := net.ResolveTCPAddr("tcp4", "127.0.0.1"+config.DefaultConfig.Network.Port)
	if err != nil {
		log.Error().Err(err).Msg("[Client] >>> Create address failed")
		return
	}
	log.Info().Msgf("[Client] >>> The address is: %s", *address)
	tcpConn, err := net.DialTCP("tcp4", nil, address)
	if err != nil {
		log.Error().Err(err).Msg("[Client] >>> Create tcpconn failed")
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
			log.Info().Msgf("Command will be run: %s", split)
			switch split[0] {
			case util.PUTCOMMAND:
				putMsg := &util.Put{
					Key:   split[1],
					Value: split[2],
				}
				marshal, _ := json.Marshal(putMsg)
				putSendMsg := tcpack.NewMessage(util.PUTID, uint32(len(marshal)), marshal)
				mp.Pack(putSendMsg)
				break
			case util.GETCOMMAND:
				getMsg := &util.Other{
					Data: split[1],
				}
				marshal, _ := json.Marshal(getMsg)
				getSendMsg := tcpack.NewMessage(util.GETID, uint32(len(marshal)), marshal)
				mp.Pack(getSendMsg)
				break
			case util.DELETECOMMAND:
				deleteMsg := &util.Other{
					Data: split[1],
				}
				marshal, _ := json.Marshal(deleteMsg)
				deleteSendMsg := tcpack.NewMessage(util.DELETEID, uint32(len(marshal)), marshal)
				mp.Pack(deleteSendMsg)
				break
			}
		}
	}()

	go func() {
		for {
			revMsg, _ := mp.Unpack()
			if revMsg.GetMsgId() == util.RESPONSEID {
				var revResule util.Other
				json.Unmarshal(revMsg.GetMsgData(), &revResule)
				fmt.Println(revResule.Data)
			}
		}
	}()
	select {}
}
