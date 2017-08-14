package worker

import (
	"math/rand"
	"github.com/go-pluto/benchmark/config"
	"github.com/go-pluto/benchmark/sessions"
)

// Functions

// expungeFolder generates an EXPUNGE command and removes
// all messages with a \Deleted flag from supplied folder.
func Generator(conf *config.Config, jobs chan Session, users []config.User) {

	// Assign jobs sessions.
	for j := 1; j <= conf.Settings.Sessions; j++ {

		// Randomly choose a user.
		i := rand.Intn(len(users))

		// Hand over the job to the worker.
		jobs <- Session{
			User:     users[i].Username,
			Password: users[i].Password,
			ID:       j,
			Commands: sessions.GenerateSession(conf.Session.MinLength, conf.Session.MaxLength),
		}
	}

	// Close jobs channel to stop all worker routines.
	// close(jobs)

}
