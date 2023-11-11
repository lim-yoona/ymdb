package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Wal struct {
		Store struct {
			Path string `yaml:"path"`
		} `yaml:"store"`
		Restore struct {
			Path string `yaml:"path"`
		} `yaml:"restore"`
		SegmentFileExt string `yaml:"segmentFileExt"`
	} `yaml:"wal"`
	Network struct {
		Port string `yaml:"port"`
	}
}

var DefaultConfig Config

func GetConfig() {
	file, err := ioutil.ReadFile("./config/ymDB.yaml")
	if err != nil {
		log.Panic().Err(err).Msg("[Config] >>> Read config file failed")
	}
	yaml.Unmarshal(file, &DefaultConfig)
}
