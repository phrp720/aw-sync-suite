package main

import (
	"aw-sync-agent/datamanager"
	"aw-sync-agent/util"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	awUrl, err := util.GetEnvVar("ACTIVITY_WATCH_URL", true)
	if err != nil {
		panic(err)
	}

	data, err := datamanager.ScrapeData(awUrl)
	print(data)
	print("\n")
	util.PromHealthCheck("http://localhost:9090")
	if err != nil {
		log.Fatal(err)
	}
	//buckets, err := aw.GetBuckets()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// Pretty-print the buckets
	//_, err = json.MarshalIndent(buckets, "", "  ")
	//if err != nil {
	//	log.Fatalf("Error marshalling buckets: %v", err)
	//}
	////fmt.Println(buckets["aw-watcher-afk_moonlight"])
	//
	////fmt.Println(string(bucketsJSON))
	//events, err := aw.GetEvents("aw-watcher-window_moonlight", nil, nil, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// Pretty-print the buckets
	//eventsJSON, err := json.MarshalIndent(events, "", "  ")
	//if err != nil {
	//	log.Fatalf("Error marshalling buckets: %v", err)
	//}
	//fmt.Println(string(eventsJSON))
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
