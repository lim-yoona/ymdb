package main

import (
	"flag"

	"github.com/lim-yoona/ymdb/config"
	"github.com/lim-yoona/ymdb/interact/server"
	"github.com/lim-yoona/ymdb/options"
	"github.com/lim-yoona/ymdb/raft"
	"github.com/lim-yoona/ymdb/route"
	"github.com/rs/zerolog"
)

var (
	storeFilePath   = flag.String("store_file_path", "/root/ymdb/walDir/store", "Specify files for ymdb to store")
	restoreFilePath = flag.String("restore_file_path", "/root/ymdb/walDir/restore", "Specify files for ymdb to restore")
)

// 服务端，之后可能会进行封装
func main() {
	flag.Parse()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	config.GetConfig()
	config.DefaultConfig.Wal.Store.Path = *storeFilePath
	config.DefaultConfig.Wal.Restore.Path = *restoreFilePath
	options.DefaultOption()
	dbServer := server.NewServer()
	go route.RouterStart(dbServer)
	raft.RaftServer(dbServer)
}
