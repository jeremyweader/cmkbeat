// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Cmkhost string `config:"cmkHost"`
	Query string `config:"query"`
	Columns []string `config:"columns"`
	Filter string `config:"filter"`
	Metrics bool `config:"metrics"`
}

var DefaultConfig = Config{
	Period: 30 * time.Second,
	Cmkhost: "localhost:6557",
	Query: "services",
	Columns: []string{"host_name", "description", "state", "plugin_output", "perf_data"},
	Filter: "",
	Metrics: true,
}

