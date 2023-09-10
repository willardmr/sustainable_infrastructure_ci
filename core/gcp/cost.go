package test

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type CloudCost struct {
	Provider     string
	RegionName   string
	Region       string
	InstanceType string
	Cost         float64
	Source       string
}

func loadCloudCosts(data [][]string) map[string]CloudCost {
	cloudCostMap := make(map[string]CloudCost)
	for i, line := range data {
		if i > 0 { // omit header line
			var rec CloudCost
			for column, field := range line {
				if column == 0 {
					rec.Provider = field
				} else if column == 1 {
					rec.RegionName = field
				} else if column == 2 {
					rec.Region = field
				} else if column == 3 {
					rec.InstanceType = field
				} else if column == 4 {
					rec.Cost, _ = strconv.ParseFloat(field, 64)
				} else if column == 5 {
					rec.Source = field
				}
			}
			cloudCostMap[rec.Region] = rec
		}
	}
	return cloudCostMap
}

func getCloudCosts() map[string]CloudCost {
	// open file
	f, err := os.Open("cost.csv")
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
	cloudCostMap := loadCloudCosts(data)

	return cloudCostMap
}
