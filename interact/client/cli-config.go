package client

import (
	"fmt"

	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/rs/zerolog/log"
)

type Nodes string

func (n Nodes) String() string {
	return string(n)
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

func InitNodeList(m map[string]string) *consistent.Consistent {
	cfg := consistent.Config{
		PartitionCount:    9,
		ReplicationFactor: 15,
		Load:              1.25,
		Hasher:            hasher{},
	}
	cst := consistent.New(nil, cfg)
	// m = make(map[string]string, 9)
	for i := 0; i < 9; i++ {
		c := byte('A' + i)
		str := "node" + string(c)
		m[str] = fmt.Sprintf("localhost:%d", 50051+i)
		member := Nodes(str)
		cst.Add(member)
	}
	return cst
}

func InitPartition(m map[int]map[string]bool) bool {
	// m = make(map[int]map[string]bool, 3)
	for i := 0; i < 3; i++ {
		p := make(map[string]bool, 3)
		for j := 0; j < 3; j++ {
			c := byte('A' + i*3 + j)
			str := "node" + string(c)
			p[str] = true
		}
		m[i] = p
	}
	return true
}

func ListNode(m map[string]string) {
	log.Info().Msg("List of nodes in a distributed key-value storage system:")
	for k, v := range m {
		s := fmt.Sprintf("node=%s==>>IP=%s", k, v)
		log.Info().Msg(s)
	}
}

func ListPartitions(m map[int]map[string]bool) {
	log.Info().Msg("List of partitions in a distributed key-value storage system")
	for k, v := range m {
		for k1, _ := range v {
			s := fmt.Sprintf("partition=%d==>>node=%s", k, k1)
			log.Info().Msg(s)
		}
	}
}
