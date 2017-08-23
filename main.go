package main

import (
	"flag"

	"io"
	"os"
	"time"

	"encoding/json"
	"math/rand"
	"net/http"
	_ "net/http/pprof"

	"cloud.google.com/go/storage"
	"github.com/go-pluto/benchmark/config"
	"github.com/go-pluto/benchmark/worker"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

// Functions

func main() {

	go func() {
		glog.Warning(http.ListenAndServe("127.0.0.1:6060", nil))
	}()

	// Parse the input flags.
	configFlag := flag.String("config", "test-config.toml", "Specify location of config file that describes test setup configuration.")
	userdbFlag := flag.String("userdb", "userdb.passwd", "Specify location of the user/password file.")
	flag.Parse()

	// Check that associated Google Cloud Project
	// is set as environment variable.
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		glog.Fatal("GOOGLE_CLOUD_PROJECT environment variable must be set")
	}

	// Make sure that we possess Application Default Credentials.
	appCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if appCredentials == "" {
		glog.Fatal("GOOGLE_APPLICATION_CREDENTIALS evironment variable must point to a valid Application Default Credentials file")
	}

	// Read configuration from file.
	conf, err := config.LoadConfig(*configFlag)
	if err != nil {
		glog.Fatalf("Error loading config: %v", err)
	}

	// Encode the configuration in json
	jsonConf, err := json.Marshal(conf)
	if err != nil {
		glog.Fatalf("Error encoding config in JSON: %v", err)
	}

	// Load users from userdb file.
	users, err := config.LoadUsers(*userdbFlag)
	if err != nil {
		glog.Fatalf("Error loading users from '%s' file: %v", *userdbFlag, err)
	}

	timestamp := time.Now()

	// Check results folder existence and create
	// a log file for this run.
	logFile, err := config.CreateLog(timestamp)
	if err != nil {
		glog.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()
	defer logFile.Sync()

	// Seed the random number generator.
	rand.Seed(conf.Settings.Seed)

	// Write first line with host information to GCS.
	// TODO comment
	_, err = logFile.WriteString("{\"Configuration\":")
	if err != nil {
		glog.Fatal(err)
	}

	// TODO comment
	_, err = logFile.Write(jsonConf)
	if err != nil {
		glog.Fatal(err)
	}

	// TODO comment
	_, err = logFile.WriteString(",\"Sessions\":[")
	if err != nil {
		glog.Fatal(err)
	}

	// Create the buffered channels. Channel "jobs" is for each session,
	// channel "logger" for the logged parameters (e.g. response time).
	jobs := make(chan worker.Session, 100)
	logger := make(chan []string, 100)

	// Start the worker pool.
	for w := 1; w <= conf.Settings.Threads; w++ {
		go worker.Worker(w, conf, jobs, logger)
	}

	go worker.Generator(conf, jobs, users)

	// Collect results and write them to disk.

	for a := 1; a <= conf.Settings.Sessions; a++ {

		logline := <-logger
		glog.Infof("Finished Session: %d", a)

		if a != 1 {
			// Write log line to log file.
			_, err := logFile.WriteString(",")
			if err != nil {
				glog.Fatal(err)
			}
		}

		for i := 0; i < len(logline); i++ {

			_, err := logFile.WriteString(logline[i])
			if err != nil {
				glog.Fatal(err)
			}
		}

		err = logFile.Sync()
		if err != nil {
			glog.Fatal(err)
		}
	}

	// TODO comment
	_, err = logFile.WriteString("]}")
	if err != nil {
		glog.Fatal(err)
	}

	err = logFile.Sync()
	if err != nil {
		glog.Fatal(err)
	}

	// Connect to GCS for log file uploading.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		glog.Fatal(err)
	}

	// Obtain writer that is able to upload
	// benchmark results to run-specific file.
	wc := client.Bucket("pluto-benchmark").Object(timestamp.Format("2006-01-02-15-04-05")).NewWriter(ctx)

	logFile.Seek(0, 0)

	_, err = io.Copy(wc, logFile)

	if err != nil {
		glog.Fatal(err)
	}

	err = wc.Close()
	if err != nil {
		glog.Fatal(err)
	}
}
