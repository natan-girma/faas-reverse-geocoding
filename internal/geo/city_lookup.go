package geo

import (
	"bufio"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type City struct {
	Name    string
	Country string
	Lat     float64
	Lon     float64
}

var (
	cities   []City
	loadOnce sync.Once
	loadErr  error
)

// LoadCities loads cities from the GeoNames cities500.txt file.
func LoadCities(path string) error {
	loadOnce.Do(func() {
		f, err := os.Open(path)
		if err != nil {
			loadErr = err
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fields := strings.Split(scanner.Text(), "\t")
			if len(fields) < 9 {
				continue
			}
			lat, err1 := strconv.ParseFloat(fields[4], 64)
			lon, err2 := strconv.ParseFloat(fields[5], 64)
			if err1 != nil || err2 != nil {
				continue
			}
			cities = append(cities, City{
				Name:    fields[1],
				Country: fields[8],
				Lat:     lat,
				Lon:     lon,
			})
		}
		if err := scanner.Err(); err != nil {
			loadErr = err
		}
	})
	return loadErr
}

// FindNearestCity returns the city and country nearest to the given lat/lon.
func FindNearestCity(lat, lon float64) (string, string, error) {
	if len(cities) == 0 {
		// Try to load the cities file from various possible locations
		possiblePaths := []string{
			os.Getenv("CITIES_DATA_PATH"),      // First check environment variable
			"/data/cities500.txt",              // Docker container path
			"data/cities500.txt",               // Local relative path
			"./data/cities500.txt",             // Explicit local path
			filepath.Join("..", "data", "cities500.txt"), // One directory up
		}
		
		// Filter out empty paths (like when env var is not set)
		var pathsToTry []string
		for _, p := range possiblePaths {
			if p != "" {
				pathsToTry = append(pathsToTry, p)
			}
		}
		
		// Try each path until one works
		var lastErr error
		for _, path := range pathsToTry {
			err := LoadCities(path)
			if err == nil {
				break // Successfully loaded
			}
			lastErr = err
		}
		
		// If we still have no cities, return the last error
		if len(cities) == 0 {
			if lastErr != nil {
				return "", "", lastErr
			}
			return "", "", errors.New("failed to load cities data from any location")
		}
	}
	
	minDist := math.MaxFloat64
	var nearest City
	for _, c := range cities {
		d := haversine(lat, lon, c.Lat, c.Lon)
		if d < minDist {
			minDist = d
			nearest = c
		}
	}
	if minDist == math.MaxFloat64 {
		return "", "", errors.New("no cities found")
	}
	return nearest.Name, nearest.Country, nil
}

// haversine calculates the great-circle distance between two points.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	lat1R := lat1 * math.Pi / 180
	lat2R := lat2 * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1R)*math.Cos(lat2R)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
} 