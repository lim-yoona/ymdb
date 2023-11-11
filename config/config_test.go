package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	GetConfig()
	fmt.Println(config)
}
