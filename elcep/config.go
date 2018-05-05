package main

import (
	"os"
	"bufio"
	"strings"
	"regexp"
	"log"
	"encoding/json"
)

type Config struct {
	Freq             *int
	ElasticsearchUrl *string
	Port             *int
	Path             *string
	QueriesFile      *string
	Queries          *map[string]string
}

var freq = 30
var host = "http://elasticsearch:9200"
var port = 8080
var path = "/metrics"
var matcherFile = "./conf/queries.cfg"

var defaultConfig = Config{
	&freq,
	&host,
	&port,
	&path,
	&matcherFile,
	nil,
}

func (config *Config) Print() {
	log.Println("Config:")
	log.Println("\tHost:", *config.ElasticsearchUrl)
	log.Println("\tFreq:", *config.Freq)
	log.Println("\tPort:", *config.Port)
	log.Println("\tQueriesFile:", *config.QueriesFile)
	log.Println("\tQueries:", prettyfy(*config.Queries))
}

func prettyfy(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func (config *Config) ReadQueriesConfig() map[string]string {
	prg := "ReadQueriesConfig()"

	var options map[string]string
	options = make(map[string]string)

	file, err := os.Open(*config.QueriesFile)
	if err != nil {
		log.Printf("%s: os.Open(): %s\n", prg, err)
		return options
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") == false {
			if strings.Contains(line, "=") == true {
				re, err := regexp.Compile(`([^=]+)=(.*)`)
				if err != nil {
					log.Printf("%s: regexp.Compile(): error=%s", prg, err)
					return options
				} else {
					config_option := re.FindStringSubmatch(line)[1]
					config_value := re.FindStringSubmatch(line)[2]
					options[config_option] = config_value
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("%s: scanner.Err(): %s\n", prg, err)
	}

	config.Queries = &options
	return options
}
