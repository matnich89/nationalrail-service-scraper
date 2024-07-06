package main

import (
	nr "github.com/matnich89/national-rail-client/nationalrail"
	"log"
	cmd "trainstats-scraper/cmd/server"
	"trainstats-scraper/config"
	"trainstats-scraper/redis"
	"trainstats-scraper/station"
)

func main() {

	c, err := config.Load()

	if err != nil {
		log.Fatalf("could not load config %v", err)
	}

	nrClient, err := nr.NewClient(
		nr.AccessTokenOpt(c.NationalRailApiKey),
	)

	if err != nil {
		log.Fatalf("could not create national rail client: %v", err)
	}

	redisClient, err := redis.NewRedisClient(c.TrainIdQueueName, c.RedisAddress)

	if err != nil {
		log.Fatalf("could not create redis client: %v", err)
	}

	stations, err := station.GetStations("./stations.txt")

	if err != nil {
		log.Fatalf("could not get stations: %v", err)
	}

	app := cmd.NewApp(nrClient, redisClient, 10, stations)

	app.SetupWorkers()

	app.Run()

}
