package server

import (
	"github.com/lim-yoona/tcpack"
	"github.com/lim-yoona/ymdb/config"
	"github.com/rs/zerolog/log"
	"net"
)

type Server struct {
	port    string
	MsQueue chan tcpack.Imessage
	ReQueue chan *tcpack.Message
}

func NewServer() *Server {
	return &Server{
		port:    config.DefaultConfig.Network.Port,
		MsQueue: make(chan tcpack.Imessage, 100),
		ReQueue: make(chan *tcpack.Message, 100),
	}
}

func (s *Server) Start() {
	address, err := net.ResolveTCPAddr("tcp4", s.port)
	if err != nil {
		log.Error().Err(err).Msg("[Server] >>> Create TCPAddr error")
	}
	log.Info().Msgf("[Server] >>> The address is: %s", s.port)
	listener, err := net.ListenTCP("tcp4", address)
	if err != nil {
		log.Error().Err(err).Msg("[Server] >>> Create listener failed")
	}
	log.Info().Msg("[Server] >>> ymDB server started!")
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Error().Err(err).Msg("[Server] >>> Accept conn failed")
		}
		mp := tcpack.NewMsgPack(8, tcpConn)
		go func() {
			for {
				msg, err := mp.Unpack()
				if err != nil {
					log.Error().Err(err).Msg("[Server] >>> Unpack on conn occur error")
				}
				s.MsQueue <- msg
			}
		}()
		go func() {
			for {
				responseMsg := <-s.ReQueue
				mp.Pack(responseMsg)
				log.Info().Msgf("[Client] >>> responseMsg has been responsed: ", responseMsg)
			}
		}()
	}
}
