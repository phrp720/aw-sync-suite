package main

import (
	"aw-sync-agent/aw"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	prometheusURL := os.Getenv("PROMETHEUS_WRITE_URL")
	if prometheusURL == "" {
		fmt.Println("Environment variable PROMETHEUS_WRITE_URL is not set or is empty")
		os.Exit(1)
	}
	buckets, err := aw.GetBuckets()
	if err != nil {
		log.Fatal(err)
	}
	// Pretty-print the buckets
	_, err = json.MarshalIndent(buckets, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling buckets: %v", err)
	}
	//fmt.Println(buckets["aw-watcher-afk_moonlight"])

	//fmt.Println(string(bucketsJSON))
	events, err := aw.GetEvents("aw-watcher-window_moonlight", nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Pretty-print the buckets
	eventsJSON, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling buckets: %v", err)
	}
	fmt.Println(string(eventsJSON))
	//client := promwrite.NewClient(prometheusURL)
	//num := float64(1 + rand.Intn(999-1))
	//fmt.Print(num)
	//_, err = client.Write(context.Background(), &promwrite.WriteRequest{
	//	TimeSeries: []promwrite.TimeSeries{
	//		// One Record
	//		{
	//			Labels: []promwrite.Label{
	//				{
	//					Name:  "__name__",
	//					Value: "ActivityWatch",
	//				},
	//			},
	//			Sample: promwrite.Sample{
	//				Time:  time.Now(),
	//				Value: 321,
	//			},
	//		},
	//		// Another Record
	//		{
	//			Labels: []promwrite.Label{
	//				{
	//					Name:  "__name__",
	//					Value: "testname",
	//				},
	//			},
	//			Sample: promwrite.Sample{
	//				Time:  time.Now(),
	//				Value: num,
	//			},
	//		},
	//	},
	//})
	//if err != nil {
	//	fmt.Printf("Failed to push: %v\n", err)
	//	return
	//}
	//fmt.Print("Push was successful!\n")
}
