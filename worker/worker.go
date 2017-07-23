package worker

import (
	"bufio"
	"fmt"
	"log"
	"time"
	s "strings"

	"crypto/tls"

	"github.com/go-pluto/benchmark/config"
	"github.com/go-pluto/benchmark/sessions"
	"github.com/golang/glog"
)

// Structs

// Session contains the user's credentials, an identifier for the
// session and a sequence of IMAP commands that has been generated
// by the sessions package.
type Session struct {
	User     string
	Password string
	ID       int
	Commands []sessions.IMAPCommand
}

// Functions

// Worker is the routine that sends the commands of the session
// to the server. The output will be logged and written in
// the logger channel.
func Worker(id int, config *config.Config, jobs chan Session, logger chan<- []string) {

	for job := range jobs {

		var output []string

		output = append(output, fmt.Sprintf("{\"SessionID\":%d,", job.ID))
		output = append(output, fmt.Sprintf("\"User\":\"%s\",", job.User))
		output = append(output, fmt.Sprintf("\"Password\":\"%s\",", job.Password))
		output = append(output, "\"Commands\":[")

		// Connect to remote server.
		tlsConn, err := tls.Dial("tcp", config.Server.Addr, &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Fatalf("Unable to connect to remote server %s: %v", config.Server.Addr, err)
		}

		conn := &Conn{
			c: tlsConn,
			r: bufio.NewReader(tlsConn),
		}

		// Login user for following IMAP commands session.
		conn.login(job.User, job.Password, id)
		glog.V(2).Info("LOGIN successful, user: ", job.User, " pw: ", job.Password)

		var commandlog []string

		for i := 0; i < len(job.Commands); i++ {

			glog.V(2).Info("Sending ", job.Commands[i].Command)

			nanos := time.Now().UnixNano()

			switch job.Commands[i].Command {

			case "CREATE":

				command := fmt.Sprintf("%dX%d CREATE %dX%s", id, i, id, job.Commands[i].Arguments[0])

				respTime, err := conn.sendSimpleCommand(command)
				if err != nil {
					log.Fatal(err)
				}

				commandlog = append(commandlog, fmt.Sprintf("[%d,\"CREATE\",%d]", nanos, respTime))

			case "DELETE":

				command := fmt.Sprintf("%dX%d DELETE %dX%s", id, i, id, job.Commands[i].Arguments[0])

				respTime, err := conn.sendSimpleCommand(command)
				if err != nil {
					log.Fatal(err)
				}

				commandlog = append(commandlog, fmt.Sprintf("[%d,\"DELETE\",%d]", nanos, respTime))

			case "APPEND":

				// command := fmt.Sprintf("%dX%d APPEND %dX%s %s %s", id, i, id, job.Commands[i].Arguments[0], job.Commands[i].Arguments[1], job.Commands[i].Arguments[2])
				command := fmt.Sprintf("%dX%d APPEND %dX%s %s", id, i, id, job.Commands[i].Arguments[0], job.Commands[i].Arguments[2])

				respTime, err := conn.sendAppendCommand(command, job.Commands[i].Arguments[3])
				if err != nil {
					log.Fatal(err)
				}

				commandlog = append(commandlog, fmt.Sprintf("[%d,\"APPEND\",%d]", nanos, respTime))

			case "SELECT":

				var command string

				if job.Commands[i].Arguments[0] == "INBOX" {
					command = fmt.Sprintf("%dX%d SELECT %s", id, i, job.Commands[i].Arguments[0])
				} else {
					command = fmt.Sprintf("%dX%d SELECT %dX%s", id, i, id, job.Commands[i].Arguments[0])
				}

				respTime, err := conn.sendSimpleCommand(command)
				if err != nil {
					log.Fatal(err)
				}

				commandlog = append(commandlog, fmt.Sprintf("[%d,\"SELECT\",%d]", nanos, respTime))

			case "STORE":

				command := fmt.Sprintf("%dX%d STORE %s FLAGS %s", id, i, job.Commands[i].Arguments[0], job.Commands[i].Arguments[1])

				respTime, err := conn.sendSimpleCommand(command)
				if err != nil {
					log.Fatal(err)
				}

				commandlog = append(commandlog, fmt.Sprintf("[%d,\"STORE\",%d]", nanos, respTime))

			case "EXPUNGE":

				command := fmt.Sprintf("%dX%d EXPUNGE", id, i)

				respTime, err := conn.sendSimpleCommand(command)
				if err != nil {
					log.Fatal(err)
				}

				commandlog = append(commandlog, fmt.Sprintf("[%d,\"EXPUNGE\",%d]", nanos, respTime))

			case "CLOSE":

				command := fmt.Sprintf("%dX%d CLOSE", id, i)

				respTime, err := conn.sendSimpleCommand(command)
				if err != nil {
					log.Fatal(err)
				}

				commandlog = append(commandlog, fmt.Sprintf("[%d,\"CLOSE\",%d]", nanos, respTime))
			}

			glog.V(2).Info(job.Commands[i].Command, " finished.")
		}


		output = append(output, s.Join(commandlog, ","))
		output = append(output, "]}")

		conn.logout(id)
		glog.V(2).Info("LOGOUT successful, user: ", job.User, " pw: ", job.Password)

		logger <- output
	}
}
