package main

import (
	"encoding/json"
	"fmt"
	"os"

	dbprovider "github.com/jamesk14022/barcrawler/database"
	types "github.com/jamesk14022/barcrawler/types"
)

type DistanceMatrix [][]float64
type RoutesMatrix [][]types.Route

type Geometry struct {
	Coordinates [][]float64 `json:"coordinates"`
}

type Route struct {
	Geometry Geometry `json:"geometry"`
	City     string   `json:"city"`
	Point1   string   `json:"point1"`
	Point2   string   `json:"point2"`
}

type PlaceIDs struct {
	PlaceIDs []string `json:"place_ids"`
}

type Location struct {
	ID     string   `json:"place_id"`
	Name   string   `json:"name"`
	Price  float32  `json:"price_level"`
	Rating float32  `json:"rating"`
	Types  []string `json:"types"`
	Photos []struct {
		PhotoReference string `json:"photo_reference"`
	} `json:"photos"`
	Geometry struct {
		Location struct {
			Latitude  float64 `json:"lat"`
			Longitude float64 `json:"lng"`
		}
	}
}

var locationDataDir = os.Getenv("LOCATION_DATA_DIR")

func getLocationData(location string) ([]types.Place, DistanceMatrix, RoutesMatrix, error) {
	var enrichedData []types.Place
	var D DistanceMatrix
	var R RoutesMatrix

	file, err := os.ReadFile("/Users/james/Documents/code/barcrawler/scrape_test/paris/info.json")
	if err != nil {
		fmt.Println("Error reading file", err)
	}
	json.Unmarshal(file, &enrichedData)

	file, err = os.ReadFile("/Users/james/Documents/code/barcrawler/scrape_test/paris/D.json")
	if err != nil {
		fmt.Println("Error reading file", err)
	}
	json.Unmarshal(file, &D)

	file, err = os.ReadFile("/Users/james/Documents/code/barcrawler/scrape_test/paris/R.json")
	if err != nil {
		fmt.Println("Error reading file", err)
	}
	json.Unmarshal(file, &R)

	return enrichedData, D, R, nil
}

func main() {
	locations := []string{"paris"}

	for _, location := range locations {

		enrichedData, _, R, _ := getLocationData(location)

		for _, place := range enrichedData {

			place.City = location

			fmt.Printf("Adding place %s\n", place.Name)
			dbprovider.Mgr.AddPlace(&place)
		}

		for i, row := range R {
			for j := range row {

				R[i][j].City = location
				R[i][j].Point1 = enrichedData[i].PlaceID
				R[i][j].Point2 = enrichedData[j].PlaceID
				dbprovider.Mgr.AddRoute(&R[i][j])
			}
		}
	}
}
