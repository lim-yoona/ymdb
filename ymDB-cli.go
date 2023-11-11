package main

import (
	"github/lim-yoona/tinyKVStore/config"
	"github/lim-yoona/tinyKVStore/interact/client"
)

func main() {
	config.GetConfig()
	client.ClientStart()
}
