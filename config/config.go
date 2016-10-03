// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Cmkhost string `config:"cmkHost"`
}

var DefaultConfig = Config{
	Period: 30 * time.Second,
	Cmkhost: "192.168.0.19:6557",
}

