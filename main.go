package main

import (
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

	if err != nil {
		log.Fatalf("could not create national rail client: %v", err)
	}

	app := cmd.NewApp(nrClient, 10, internal.GetStations())

	app.SetupWorkers()

	app.Run()

}
