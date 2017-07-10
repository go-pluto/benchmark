package main

import (
	"flag"

	"math/rand"

	"github.com/go-pluto/benchmark/config"
	"github.com/go-pluto/benchmark/sessions"
	"github.com/go-pluto/benchmark/worker"
	"github.com/golang/glog"
)

// Functions

func main() {

	// Parse the input flags.
	configFlag := flag.String("config", "test-config.toml", "Specify location of config file that describes test setup configuration.")
	userdbFlag := flag.String("userdb", "userdb.passwd", "Specify location of the user/password file.")
	flag.Parse()

	// Read configuration from file.
	conf, err := config.LoadConfig(*configFlag)
	if err != nil {
		glog.Fatalf("Error loading config: %v", err)
	}

	// Load users from userdb file.
	users, err := config.LoadUsers(*userdbFlag)
	if err != nil {
		glog.Fatalf("Error loading users from '%s' file: %v", *userdbFlag, err)
	}

	// Check results folder existence and create
	// a log file for this run.
	logFile, err := config.CreateLog()
	if err != nil {
		glog.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()
	defer logFile.Sync()

	// Seed the random number generator.
	rand.Seed(conf.Settings.Seed)

	// Create the buffered channels. Channel "jobs" is for each session,
	// channel "logger" for the logged parameters (e.g. response time).
	jobs := make(chan worker.Session, 100)
	logger := make(chan []string, 100)

	// Start the worker pool.
	for w := 1; w <= conf.Settings.Threads; w++ {
		go worker.Worker(w, conf, jobs, logger)
	}

	// Assign jobs sessions.
	for j := 1; j <= conf.Settings.Sessions; j++ {

		// Randomly choose a user.
		i := rand.Intn(len(users))

		// Hand over the job to the worker.
		jobs <- worker.Session{
			User:     users[i].Username,
			Password: users[i].Password,
			ID:       j,
			Commands: sessions.GenerateSession(conf.Session.MinLength, conf.Session.MaxLength),
		}
	}

	glog.Infof("Generated %d sessions.", conf.Settings.Sessions)

	// Close jobs channel to stop all worker routines.
	close(jobs)

	// Collect results and write them to disk.
	for a := 1; a <= conf.Settings.Sessions; a++ {

		logline := <-logger

		for i := 0; i < len(logline); i++ {
			logFile.WriteString(logline[i])
			glog.Infof("%s", logline[i])
			logFile.WriteString("\n")
		}
	}
}
