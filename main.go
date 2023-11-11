package main

import (
	"github.com/rs/zerolog"
	"github/lim-yoona/tinyKVStore/config"
	"github/lim-yoona/tinyKVStore/interact/server"
	"github/lim-yoona/tinyKVStore/options"
	"github/lim-yoona/tinyKVStore/route"
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
