package main

import (
	"aw-sync-agent/synchronizer"
	"aw-sync-agent/util"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	//Here we must load the env variables and check the flags
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	awUrl, err := util.GetEnvVar("ACTIVITY_WATCH_URL", true)
	if err != nil {
		panic(err) // force quit. mandatory url
	}
	prometheusUrl, err := util.GetEnvVar("PROMETHEUS_URL", true)
	if err != nil {
		panic(err) // force quit. mandatory url
	}
	err = synchronizer.Start(awUrl, prometheusUrl)
	if err != nil {
		panic(err) // handle if something wrong happens
	}
}
