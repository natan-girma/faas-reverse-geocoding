package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saidsef/faas-reverse-geocoding/internal/geo"
)

type ReverseGeocodeResponse struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// ReverseGeocodeHandler handles /reverse-geocode?lat=...&lon=...
func ReverseGeocodeHandler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	if latStr == "" || lonStr == "" {
		http.Error(w, "lat and lon are required", http.StatusBadRequest)
		return
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid lat", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid lon", http.StatusBadRequest)
		return
	}
	city, country, err := geo.FindNearestCity(lat, lon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := ReverseGeocodeResponse{City: city, Country: country}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
} 