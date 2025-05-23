package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saidsef/faas-reverse-geocoding/internal/handlers"
	"github.com/saidsef/faas-reverse-geocoding/internal/metrics"
	"github.com/saidsef/faas-reverse-geocoding/internal/utils"
)

var (
	// port defines the port number on which the server will listen.
	port int

	// verbose is a flag that indicates whether verbose logging is enabled.
	verbose = utils.Verbose

	// cache defines the cache duration in minutes.
	cache = handlers.CACHE_DURATION_MINUTES
)

// loggingMiddleware is an HTTP middleware that logs the details of each incoming request.
// It logs the remote address, HTTP method, URL, content length, host, and protocol of the request.
//
// Parameters:
// - next: The next http.HandlerFunc to be called after logging the request details.
//
// Returns:
// - An http.HandlerFunc that logs the request details and then calls the next handler.
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Logger.Infof("%s %s %s %d %s %s", r.RemoteAddr, r.Method, r.URL, r.ContentLength, r.Host, r.Proto)
		next.ServeHTTP(w, r)
	}
}

func main() {
	flag.IntVar(&port, "port", 8080, "HTTP listening PORT")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.IntVar(&cache, "cache", 60, "Cache duration minutes (default 30)")
	flag.Parse()

	// Set the verbosity level in the utils package
	utils.SetVerbose(verbose)

	// Initialise registers of the metrics Hostname
	metrics.Init()

	// Set cache duration minutes
	handlers.SetCacheDurationMinutes(cache)

	r := http.NewServeMux()
	r.HandleFunc("/", loggingMiddleware(handlers.LatitudeLongitude))
	r.HandleFunc("/reverse-geocode", loggingMiddleware(handlers.ReverseGeocodeHandler))
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           r,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	utils.Logger.Infof("Server is running on port %d and address %s", port, srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		utils.Logger.Fatal(err)
	}
}
