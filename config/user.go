package config

import (
	"bytes"

	"io/ioutil"
)

// Structs

// User represents a user by name and password.
type User struct {
	Username string
	Password string
}

// Functions

// LoadUsers populates all users from supplied
// file into a slice of above User struct.
func LoadUsers(userdbFile string) ([]User, error) {

	users := make([]User, 0, 30)

	// Load whole file content.
	content, err := ioutil.ReadFile(userdbFile)
	if err != nil {
		return nil, err
	}

	// Split content at newline.
	lines := bytes.Split(content, []byte("\n"))

	for _, line := range lines {

		if len(line) > 0 {

			// Split line at ':{plain}'.
			userData := bytes.Split(line, []byte(":{plain}"))

			// Append new User element with parsed data.
			users = append(users, User{
				Username: string(userData[0]),
				Password: string(userData[1]),
			})
		}
	}

	return users, nil
}
