// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Cmkhost string `config:"cmkHost"`
	Query string `config:"query"`
	Columns []string `config:"columns"`
	Filter []string `config:"filter"`
	Metrics bool `config:"metrics"`
}

var DefaultConfig = Config{
	Period: 30 * time.Second,
	Cmkhost: "localhost:6557",
	Query: "services",
	Columns: []string{"host_name", "host_alias", "display_name", "check_command", "state", "acknowledged", "plugin_output", "long_plugin_output", "percent_state_change", "perf_data"},
	Filter: []string{"checks_enabled = 1", "acknowledged = 0"},
	Metrics: true,
}

