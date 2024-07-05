package main

import (
	"github.com/go-redis/redis/v8"
	nr "github.com/matnich89/national-rail-client/nationalrail"
	"log"
	"os"
	cmd "trainstats-scraper/cmd/server"
	"trainstats-scraper/internal"
)

func main() {
	nrClient, err := nr.NewClient(
		nr.AccessTokenOpt(os.Getenv("NATIONAL_RAIL_API_KEY")),
	)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err != nil {
		log.Fatalf("could not create national rail client: %v", err)
	}

	stations, err := internal.GetStations("./stations.txt")

	if err != nil {
		log.Fatalf("could not get stations: %v", err)
	}

	app := cmd.NewApp(nrClient, redisClient, 10, stations)

	app.SetupWorkers()

	app.Run()

}
