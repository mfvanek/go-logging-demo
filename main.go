package main

import (
	"encoding/json"
	"fmt"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"time"
)

var logger = logrus.New()
var log = logger.WithFields(logrus.Fields{
	"appType":           "GO",
	"envType":           "K8S",
	"appName":           "go-logging-demo-app",
	"agrType":           "OPENSHIFT_EVENT",
	"osEventActionType": "ADD",
})
var router = mux.NewRouter().StrictSlash(true)

func logToJson(w http.ResponseWriter, r *http.Request) {
	extEventId := r.URL.Query().Get("extEventId")

	if extEventId == "" {
		extEventId = uuid.New().String()
	}

	log.WithFields(logrus.Fields{
		"eventId": uuid.New().String(),
	}).Info("test log! extEventId = ", extEventId)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(time.Now())
}

func livenessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleRequests() {
	router.Path("/v1/log").Queries("extEventId", "{extEventId}").HandlerFunc(logToJson)
	router.Path("/v1/log").HandlerFunc(logToJson)
	router.Path("/metrics").Handler(promhttp.Handler())
	router.HandleFunc("/liveness", livenessHandler)
	router.HandleFunc("/readiness", readinessHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	host, hostExists := os.LookupEnv("LOG_COLLECTOR_HOST")
	if !hostExists {
		host = "127.0.0.1"
	}
	tcpPort, tcpPortExists := os.LookupEnv("LOG_COLLECTOR_PORT")
	if !tcpPortExists {
		tcpPort = "5170"
	}
	_, _ = fmt.Fprintf(os.Stdout, "Logs collector host %s port %s\n", host, tcpPort)

	conn, err := net.Dial("tcp", host+":"+tcpPort)
	checkError(err)

	if err == nil {
		hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{}))
		logger.Hooks.Add(hook)
	}

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	handleRequests()
}

func checkError(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Logs collector connection error: %s", err.Error())
	}
}
