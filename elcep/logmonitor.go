package main

import (
	"strings"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"log"
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

const query_path = "/_search"
const queryTemplate = `{
    "query": {
        "query_string": {
            "query": "<query>"
        }
    },
    "size":0
}`

type LogMonitor struct {
	Name string

	request struct {
		body string
	}
	response ElasticResponse

	query     string
	url       string
	LastCount *float64

	metrics struct {
		matchCounter         prometheus.Counter
		rpcDurationHistogram prometheus.Histogram
	}
}

func (logMon *LogMonitor) Build(config Config, query string) {
	logMon.LastCount = new(float64)
	logMon.url = fmt.Sprintf(*config.ElasticsearchUrl + query_path)

	logMon.buildQuery(query)

	logMon.metrics.matchCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "logs_matched_" + logMon.Name + "_total",
		Help: "Counts number of matched logs for " + logMon.Name,
	})
	logMon.metrics.rpcDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "logs_matched_rpc_durations_" + logMon.Name + "_histogram_seconds",
		Help:    "Logs matched RPC latency distributions for " + logMon.Name,
		Buckets: prometheus.LinearBuckets(0.001, 0.001, 10),
	})
}

func (logMon *LogMonitor) Register() {
	prometheus.MustRegister(logMon.metrics.matchCounter)
	prometheus.MustRegister(logMon.metrics.rpcDurationHistogram)
}

func (logMon *LogMonitor) Perform() {
	increment := logMon.countLogs()

	if increment < 0 {
		increment = 0
	}
	logMon.metrics.matchCounter.Add(increment)
}

func (logMon *LogMonitor) countLogs() float64 {
	start := time.Now()
	response := logMon.execRequest()
	end := time.Now()

	duration := end.Sub(start).Seconds()
	logMon.metrics.rpcDurationHistogram.Observe(duration)

	increment := response.HitCount() - *logMon.LastCount
	*logMon.LastCount = response.HitCount()

	return increment
}

func (logMon *LogMonitor) buildQuery(query string) {
	logMon.query = query
	logMon.request.body = strings.Replace(queryTemplate, "<query>", query, 1)
}

func (logMon *LogMonitor) execRequest() ElasticResponse {
	req, err := http.NewRequest("GET", logMon.url, bytes.NewBufferString(logMon.request.body))
	req.Header.Set("Content-Type", "application/Json")

	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Do: ", err)
		return logMon.response
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&logMon.response); err != nil {
		log.Println(err)
	}

	return logMon.response
}
