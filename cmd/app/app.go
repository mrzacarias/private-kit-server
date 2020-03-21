package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mrzacarias/private-kit-server/config"
	"github.com/mrzacarias/private-kit-server/database"
	infected "github.com/mrzacarias/private-kit-server/internal/infected"
	"github.com/mrzacarias/private-kit-server/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// DBClient will handle the Database connections
var DBClient *sql.DB

// InfectedClient is a (mockable) Infected client
var InfectedClient infected.Contract

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	DBClient = database.NewDBClient()
	InfectedClient = infected.NewClient(DBClient)
}

// HealthCheckHandler is private-kit-server endpoint for livenessProbe
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// InfectedHandler is private-kit-server endpoint for the internal package Infected
func InfectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	metrics.RequestsTotal.With(prometheus.Labels{"endpoint": "infected"}).Inc()

	// Total request time start
	requestStart := time.Now()

	var req infected.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	switch {
	case err == io.EOF:
		metrics.RequestsErrors.With(prometheus.Labels{"endpoint": "infected", "type": "json.NewDecoder(r.Body).Decode(&req)"}).Inc()
		log.WithFields(log.Fields{"endpoint": "infected"}).Errorln("POST body empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	case err != nil:
		metrics.RequestsErrors.With(prometheus.Labels{"endpoint": "infected", "type": "json.NewDecoder(r.Body).Decode(&req)"}).Inc()
		log.WithFields(log.Fields{"endpoint": "infected"}).Errorln("Error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// Make infected.StoreInfected request
	err = InfectedClient.StoreInfected(req)
	if err != nil {
		metrics.RequestsErrors.With(prometheus.Labels{"endpoint": "infected", "type": "infected.StoreInfected"}).Inc()
		log.WithFields(log.Fields{"endpoint": "infected", "error": err}).Errorln("Error on infected.StoreInfected")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metrics.RequestDurationTotal.With(prometheus.Labels{"endpoint": "infected"}).Set(time.Since(requestStart).Seconds())
	w.WriteHeader(http.StatusOK)
}

// serveHTTP will start the HTTP server and it's endpoints
func serveHTTP(srv *http.Server, errChan chan error) {
	http.HandleFunc("/", HealthCheckHandler)
	http.HandleFunc("/healthcheck", HealthCheckHandler)
	http.HandleFunc("/healthz", HealthCheckHandler)
	http.HandleFunc("/infected", InfectedHandler)

	log.WithFields(log.Fields{"address": srv.Addr}).Info("`private-kit-server` listening")
	err := srv.ListenAndServe()
	if err != nil {
		errChan <- fmt.Errorf("cannot listen to address: %v", err)
		return
	}
}

// serveMetrics will start the HTTP server for metrics, that will be consumed by telegraf
func serveMetrics(srv *http.Server, errChan chan error) {
	mux := http.NewServeMux()
	srv.Handler = mux
	mux.Handle("/metrics", promhttp.Handler())

	log.WithFields(log.Fields{"address": srv.Addr}).Info("`private-kit-server` Metrics listening")
	err := srv.ListenAndServe()
	if err != nil {
		errChan <- fmt.Errorf("cannot listen to address: %v", err)
		return
	}
}

// Main thread
func main() {
	log.Info("private-kit-server Initialized!")
	cfg := config.GetConfig()

	// Preparing channels to listen to hard errors and signals
	errChan := make(chan error)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	// HTTP logic
	proxySrv := &http.Server{Addr: fmt.Sprintf(":%s", cfg.Port)}
	metricsSrv := &http.Server{Addr: fmt.Sprintf(":%s", cfg.MetricsPort)}
	go serveHTTP(proxySrv, errChan)
	go serveMetrics(metricsSrv, errChan)

	// Deferring the database connection closing
	defer DBClient.Close()

	// Blocking the main thread execution while waits for a hard error or signal
	select {
	case err := <-errChan:
		log.WithFields(log.Fields{"error": err}).Error("Failed to start server")
	case sig := <-sigChan:
		if sigMsg := sig.String(); sigMsg == "interrupt" || sigMsg == "terminated" {
			log.WithFields(log.Fields{"signal_message": sigMsg}).Error("Received termination signal, shutting down servers...")
		}
	}
}
