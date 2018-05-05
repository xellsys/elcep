package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
	"strconv"
)

func main() {
	var config Config
	var logMons []LogMonitor

	config.ElasticsearchUrl = flag.String("url", *defaultConfig.ElasticsearchUrl, "The elastic search endpoint")
	config.Port = flag.Int("port", *defaultConfig.Port, "The port to listen on for HTTP requests")
	config.Path = flag.String("path", *defaultConfig.Path, "The path to listen on for HTTP requests")
	config.Freq = flag.Int("freq", *defaultConfig.Freq, "The interval in seconds in which to query elastic search")
	config.QueriesFile =
		flag.String("config", *defaultConfig.QueriesFile, "The path to the queries.cfg")
	flag.Parse()

	config.ReadQueriesConfig()
	config.Print()

	logMons = make([]LogMonitor, len(*config.Queries))
	i := 0
	for name, query := range *config.Queries {
		logMons[i].Name = name
		logMons[i].Build(config, query)
		logMons[i].Register()
		i += 1
	}

	// logMons loop
	go func() {
		for {
			for _, logMon := range logMons {
				logMon.Perform()

			}
			time.Sleep(time.Duration(*config.Freq) * time.Second)
		}
	}()

	http.Handle(*config.Path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*config.Port), nil))
}
