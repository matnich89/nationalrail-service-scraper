package station

import (
	"github.com/matnich89/national-rail-client/nationalrail"
	"os"
	"strings"
)

func GetStations(filename string) ([]nationalrail.CRSCode, error) {

	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	var stations []nationalrail.CRSCode

	for _, line := range lines {
		code := strings.Trim(strings.TrimSpace(line), "\"")

		if code != "" {
			stations = append(stations, nationalrail.CRSCode(code))
		}
	}

	return stations, nil
}
