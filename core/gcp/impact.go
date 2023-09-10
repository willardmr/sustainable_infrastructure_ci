package test

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

type CloudIntensity struct {
	Provider      string
	Region        string
	Location      string
	Impact        float64
	Source        string
	GeneralRegion string
}

func loadCloudIntensities(data [][]string) []CloudIntensity {
	var cloudIntensityList []CloudIntensity
	for i, line := range data {
		if i > 0 { // omit header line
			var rec CloudIntensity
			for column, field := range line {
				if column == 0 {
					rec.Provider = field
				} else if column == 1 {
					rec.Region = field
					switch rec.Provider {
					case "gcp":
						rec.GeneralRegion = getGcpGeneralRegion(field)
					case "aws":
						rec.GeneralRegion = getAwsGeneralRegion(field)
					case "azure":
						rec.GeneralRegion = getAzureGeneralRegion(field)
					}
				} else if column == 2 {
					rec.Location = field
				} else if column == 3 {

					rec.Impact, _ = strconv.ParseFloat(field, 64)
				} else if column == 4 {
					rec.Source = field
				}
			}
			cloudIntensityList = append(cloudIntensityList, rec)
		}
	}
	return cloudIntensityList
}

func getGcpGeneralRegion(region string) string {
	switch strings.Count(region, "-") {
	case 0:
		return region
	case 1:
		// Takes the form of us-west2
		return region[:len(region)-1]
	case 2:
		// Takes the form of us-west2-b
		topLevelRegion := strings.Join(strings.Split(region, "-")[:2], "-")
		return topLevelRegion[:len(topLevelRegion)-1]
	default:
		return region
	}
}

func getAwsGeneralRegion(region string) string {
	return region[:len(region)-2]
}

func getAzureGeneralRegion(region string) string {
	// Azure regions do not follow a consistent pattern, ignore for now
	return region
}

func getCloudIntensities() []CloudIntensity {
	// open file
	f, err := os.Open("impact.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	cloudIntensityList := loadCloudIntensities(data)

	return cloudIntensityList
}
