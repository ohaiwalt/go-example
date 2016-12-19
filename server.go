package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// production-readiness standards
// 1. /healthz endpoint to report health of service
// 2. /metrics endpoint served by prometheus library of choice
// 3. Logging events to stdout
// 4. Confluence documentation of service, example as README.md here

var requestCount = prom.NewCounter(prom.CounterOpts{
	Name: "http_request_total",
	Help: "HTTP Request Count.",
})

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World, %s!", req.URL.Path[1:])
	log.WithFields(log.Fields{
		"path": req.URL.Path,
	}).Info("Request")
	requestCount.Inc()
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, http.StatusOK)
}

func init() {
	prom.MustRegister(requestCount)
}

func main() {
	logger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})
	logger.Info("Starting app...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", healthzHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
