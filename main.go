package main

import (
	"aw-sync-agent/datamanager"
	"aw-sync-agent/prometheus"
	"aw-sync-agent/util"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"strings"
)

func main() {
	//Here we must load the env variables and check the flags
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	awUrl, err := util.GetEnvVar("ACTIVITY_WATCH_URL", true)
	if err != nil {
		panic(err)
	}
	prometheusUrl, err := util.GetEnvVar("PROMETHEUS_URL", true)
	if err != nil {
		panic(err)
	}
	prometheusClient := prometheus.NewClient(fmt.Sprintf("%s%s", prometheusUrl, "/api/v1/write"))
	scrapedData, err := datamanager.ScrapeData(awUrl)
	for watcher, data := range scrapedData {
		aggregatedData := datamanager.AggregateData(data, strings.ReplaceAll(watcher, "-", "_")) //metric names must not have '-'
		err = datamanager.PushData(prometheusClient, prometheusUrl, aggregatedData)
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}
}
