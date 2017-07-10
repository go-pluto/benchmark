package worker

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"crypto/tls"

	"github.com/golang/glog"
)

// Structs

// Conn encapsulates connection adapters to write
// and read from an active TLS connection.
type Conn struct {
	c *tls.Conn
	r *bufio.Reader
}

// Functions

// login sends a LOGIN command with the given
// username/password combination on given TLS
// connection.
func (c *Conn) login(username string, password string, id int) error {

	okAnswer := fmt.Sprintf("%dX OK", id)

	// Consume mandatory IMAP greeting.
	_, err := c.r.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error during receiving initial server greeting: %v", err)
	}

	// Send LOGIN command with parameters.
	_, err = fmt.Fprintf(c.c, "%dX LOGIN %s %s\r\n", id, username, password)
	if err != nil {
		return fmt.Errorf("sending LOGIN to server failed with: %v", err)
	}

	// Wait for success message.
	answer, err := c.r.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error receiving answer to LOGIN as user: %v", err)
	}

	// Check for success indicator in answer.
	for !strings.Contains(answer, okAnswer) {

		nextAnswer, err := c.r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error during receiving nextAnswer: %v", err)
		}

		answer = nextAnswer
	}

	return nil
}

// sendSimpleCommand sends an IMAP command string
// on given connection. The time between sending
// the message and receiving the corresponding
// confirmation will be measured and returned.
func (c *Conn) sendSimpleCommand(command string) (int64, error) {

	okAnswer := strings.Split(command, " ")[0]

	glog.V(3).Info("Sending command: ", command)

	// Start time taken here.
	timeStart := time.Now().UnixNano()

	_, err := fmt.Fprintf(c.c, "%s\r\n", command)
	if err != nil {
		return -1, fmt.Errorf("error during sending: %v", err)
	}

	answer, err := c.r.ReadString('\n')
	if err != nil {
		return -1, fmt.Errorf("error during receiving after first imap command: %v", err)
	}

	glog.V(3).Info("Answer: ", answer)

	for !strings.HasPrefix(answer, okAnswer) {

		nextAnswer, err := c.r.ReadString('\n')
		if err != nil {
			return -1, fmt.Errorf("error during receiving after other imap command: %v", err)
		}

		answer = nextAnswer
		glog.V(3).Info("Answer: ", answer)
	}

	// End time taken here.
	timeEnd := time.Now().UnixNano()

	if !strings.Contains(answer, "OK") {
		glog.Warningf("server responded unexpectedly to command: %s\n by answer: %s", command, answer)
	}

	return (timeEnd - timeStart), nil
}

// sendAppendCommand sends an IMAP command string
// that contains an APPEND command on the given
// connection "con". The time between the send of
// the message and the receive of the imap confirmation
// will be counted and returned.
func (c *Conn) sendAppendCommand(command string, literal string) (int64, error) {

	glog.V(3).Info("Sending command: ", command)

	okAnswer := strings.Split(command, " ")[0]

	// Start time taken here.
	timeStart := time.Now().UnixNano()

	_, err := fmt.Fprintf(c.c, "%s\r\n", command)
	if err != nil {
		return -1, fmt.Errorf("error during sending: %v", err)
	}

	answer, err := c.r.ReadString('\n')
	if err != nil {
		return -1, fmt.Errorf("error during receiving after append command: %v", err)
	}

	glog.V(3).Info("Answer: ", answer)

	if (answer != "+ OK\r\n") && (answer != "+ Ready for literal data\r\n") {
		return -1, fmt.Errorf("did not receive continuation command from server")
	}

	// Send message literal.
	_, err = fmt.Fprintf(c.c, "%s\r\n", literal)
	if err != nil {
		return -1, fmt.Errorf("sending mail message to server failed with: %v", err)
	}

	answer, err = c.r.ReadString('\n')
	if err != nil {
		return -1, fmt.Errorf("error during receiving response to APPEND: %v", err)
	}

	glog.V(3).Info("Answer: ", answer)

	for !strings.HasPrefix(answer, okAnswer) {

		nextAnswer, err := c.r.ReadString('\n')
		if err != nil {
			return -1, fmt.Errorf("error during receiving after append literal: %v", err)
		}

		answer = nextAnswer
		glog.V(3).Info("Answer: ", answer)
	}

	// End time taken here.
	timeEnd := time.Now().UnixNano()

	if !strings.Contains(answer, "OK") {
		glog.Warningf("server responded unexpectedly to command: %s", command)
	}

	return (timeEnd - timeStart), nil
}

// logout sends a LOGOUT command to the server.
func (c *Conn) logout(id int) error {

	okAnswer := fmt.Sprintf("%dZ", id)

	_, err := fmt.Fprintf(c.c, "%dZ LOGOUT\r\n", id)
	if err != nil {
		return fmt.Errorf("error during LOGOUT: %v", err)
	}

	answer, err := c.r.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error receiving first part of LOGOUT response: %v", err)
	}

	for !strings.Contains(answer, okAnswer) {

		nextAnswer, err := c.r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error during LOGOUT: %v", err)
		}

		answer = nextAnswer
	}

	return nil
}
