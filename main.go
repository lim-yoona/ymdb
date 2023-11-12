package main

import (
	"github.com/lim-yoona/ymdb/config"
	"github.com/lim-yoona/ymdb/interact/server"
	"github.com/lim-yoona/ymdb/options"
	"github.com/lim-yoona/ymdb/route"
	"github.com/rs/zerolog"
)

// 服务端，之后可能会进行封装
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	config.GetConfig()
	options.DefaultOption()
	dbServer := server.NewServer()
	go route.RouterStart(dbServer)
	dbServer.Start()
}
