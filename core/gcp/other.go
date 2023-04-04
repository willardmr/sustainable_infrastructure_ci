package test 

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

type CloudIntensity struct {
	Provider       string
	ProviderName   string
	OffsetRatio    int
	Region         string
	RegionName     string
	CountryName    string
	CountryIsoCode string
	State          string
	City           string
	Impact         string
	Source         string
	GeneralRegion  string
}

func createShoppingList(data [][]string) []CloudIntensity {
	var cloudIntensityList []CloudIntensity
	for i, line := range data {
		if i > 0 { // omit header line
			var rec CloudIntensity
			for j, field := range line {
				if j == 0 {
					rec.Provider = field
				} else if j == 1 {
					rec.ProviderName = field
				} else if j == 2 {
					rec.OffsetRatio, _ = strconv.Atoi(field)
					//rec.OffsetRatio = offsetRatio
				} else if j == 3 {
					rec.Region = field
					switch rec.Provider {
					case "gcp":
						rec.GeneralRegion = getGcpGeneralRegion(field)
					case "aws":
						rec.GeneralRegion = getAwsGeneralRegion(field)
					case "azure":
						rec.GeneralRegion = getAzureGeneralRegion(field)

					}
				} else if j == 4 {
					rec.RegionName = field
				} else if j == 5 {
					rec.CountryName = field
				} else if j == 6 {
					rec.CountryIsoCode = field
				} else if j == 7 {
					rec.State = field
				} else if j == 8 {
					rec.City = field
				} else if j == 9 {
					rec.Impact = field
				} else if j == 10 {
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
	cloudIntensityList := createShoppingList(data)

	return cloudIntensityList
}

