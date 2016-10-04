package beater

import (
	"fmt"
	"time"
	"strings"
	"regexp"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/jeremyweader/cmkbeat/config"
	"github.com/jeremyweader/go-livestatus"
)

type Cmkbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Cmkbeat{
		done: make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Cmkbeat) Run(b *beat.Beat) error {
	
	if len(bt.config.Cmkhost) < 1 {
		return fmt.Errorf("Error: Invalid cmkhost config \"%s\"", bt.config.Cmkhost)
	}
	if len(bt.config.Query) < 1 {
		return fmt.Errorf("Error: Invalid query config \"%s\"", bt.config.Query)
	}
	if len(bt.config.Columns) < 1 {
		return fmt.Errorf("Error: Invalid columns config \"%s\"", bt.config.Columns)
	}
	
	
	logp.Info("------Config-------")
	logp.Info("Host: %s", bt.config.Cmkhost)
	logp.Info("Query: %s", bt.config.Query)
	logp.Info("Columns: %s", bt.config.Columns)
	logp.Info("Filter: %s", bt.config.Filter)
	logp.Info("Metrics: %t", bt.config.Metrics)
	logp.Info("--------------")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)

	for {
		select {
			case <-bt.done:
				return nil
			case <-ticker.C:
		}

		err := bt.lsQuery(bt.config.Cmkhost, b.Name)
		if err != nil {
			logp.Warn("Error executing query: %s", err)
			return err
		}
	}
}

func (bt *Cmkbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Cmkbeat) lsQuery(lshost string, beatname string) error {
	
	logp.Info("Starting query")
    start := time.Now()
	
	var host string = lshost
	var query string = bt.config.Query
	var columns []string = bt.config.Columns
	var filter string = bt.config.Filter
	var metrics bool = bt.config.Metrics
	
    l := livestatus.NewLivestatus("tcp", host)
    q := l.Query(query)
    q.Columns(columns)
	
	if len(filter) > 0 {
		q.Filter(filter)
	}

    resp, err := q.Exec()
    if err != nil {
		return err
    }

	numRecords := 0

    for _, r := range resp.Records {
        host_name, err := r.GetString("host_name")
		description, err := r.GetString("description")
		state, err := r.GetInt("state")
		plugin_output, err := r.GetString("plugin_output")
		perf_data, err := r.GetString("perf_data")
		if err != nil {
			logp.Warn("Problem parsing response fields: %s", err)
		}
		
		logp.Info("hostname: %s", host_name)
		logp.Info("description: %s", description)
		logp.Info("state: %v", state)
		logp.Info("plugin_output: %s", plugin_output)
		logp.Info("perfdata: %s", perf_data)
		
		event := common.MapStr {
			"@timestamp":	common.Time(time.Now()),
			"type":		beatname,
			"host":		host_name,
			"description":	description,
			"state":	state,
			"output":	plugin_output,
			"perfdata":	perf_data,
		}
		
		if metrics == true {
			if len(perf_data) > 0 {
				var perfObjMap map[string]map[string]string
				var perfDataSplit []string

				perfDataSplit = strings.Split(perf_data, " ")
				perfObjMap = make(map[string]map[string]string)
				for _, perfObj := range perfDataSplit {
					var perfObjSplit []string
					var dataSplit []string
					if perfObj != "" {
						perfObjSplit = strings.Split(perfObj, "=")
						if len(perfObjSplit) == 2 {
							item := perfObjSplit[0]
							data := perfObjSplit[1]
							if data != "" {
								if strings.Contains(data, ";") {
									dataSplit = strings.Split(data, ";")
									perfObjMap[item] = make(map[string]string)
									dsLen := len(dataSplit)
									if dsLen >= 1 {
										if len(dataSplit[0]) > 0 {
											re := regexp.MustCompile("[0-9\\.]+")
											num := re.FindAllString(dataSplit[0], 1)
											if len(num) > 0 {
												perfObjMap[item]["value"] = num[0]
											}
											logp.Info("metrics: %s: value: %v", item, num[0])
										}
									}
									if dsLen >= 2 {
										if len(dataSplit[1]) > 0 {
											re := regexp.MustCompile("[0-9\\.]+")
											num := re.FindAllString(dataSplit[1], 1)
											if len(num) > 0 {
												perfObjMap[item]["min"] = num[0]
											}
											logp.Info("metrics: %s: min: %v", item, num[0])
										}
									}
									if dsLen >= 3 {
										if len(dataSplit[2]) > 0 {
											re := regexp.MustCompile("[0-9\\.]+")
											num := re.FindAllString(dataSplit[2], 1)
											if len(num) > 0 {
												perfObjMap[item]["max"] = num[0]
											}
											logp.Info("metrics: %s: max: %v", item, num[0])
										}
									}
									if dsLen >= 4 {
										if len(dataSplit[3]) > 0 {
											re := regexp.MustCompile("[0-9\\.]+")
											num := re.FindAllString(dataSplit[3], 1)
											if len(num) > 0 {
												perfObjMap[item]["warn"] = num[0]
											}
											logp.Info("metrics: %s: warn: %v", item, num[0])
										}
									}
									if dsLen >= 5 {
										if len(dataSplit[4]) > 0 {
											re := regexp.MustCompile("[0-9\\.]+")
											num := re.FindAllString(dataSplit[4], 1)
											if len(num) > 0 {
												perfObjMap[item]["crit"] = num[0]
											}
											logp.Info("metrics: %s: crit: %v", item, num[0])
										}
									}
								} else {
									perfObjMap[item] = make(map[string]string)
									re := regexp.MustCompile("[0-9\\.]+")
									num := re.FindAllString(data, 1)
									if len(num) > 0 {
										perfObjMap[item]["value"] = num[0]
									}
									logp.Info("metrics: %s: value: %v", item, num[0])
								}
							}
						}
					}
				}
				event["metrics"] = perfObjMap
			}
		}
		logp.Info("Publishing event")
		bt.client.PublishEvent(event)
		numRecords++
    }
    elapsed := time.Since(start)
    logp.Info("%v events submitted in %s.", numRecords, elapsed)
    return nil
}

