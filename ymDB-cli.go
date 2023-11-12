package main

import (
	"github.com/lim-yoona/ymdb/config"
	"github.com/lim-yoona/ymdb/interact/client"
)

func main() {
	config.GetConfig()
	client.ClientStart()
}
