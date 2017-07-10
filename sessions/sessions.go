package sessions

import (
	"fmt"
	"strconv"

	"math/rand"

	"github.com/benchmark/utils"
)

// Structs

// IMAPCommand contains the string of the command
// and the corresponding arguments.
type IMAPCommand struct {
	Command   string
	Arguments []string
}

// Folder represents an imap folder with the
// contained messages.
type Folder struct {
	Foldername string
	Messages   []Message
}

// Message represents a message, in this case
// only the flags are relevant to generate a session.
type Message struct {
	Flags []string
}

// Functions

// expungeFolder generates an EXPUNGE command and
// removes all messages with a \Deleted flag from a folder.
func expungeFolder(folder *Folder) IMAPCommand {

	for j := 0; j < len(folder.Messages); j++ {

		for k := 0; k < len(folder.Messages[j].Flags); k++ {

			if folder.Messages[j].Flags[k] == "\\Deleted" {
				folder.Messages = append(folder.Messages[:j], folder.Messages[j+1:]...)
				j = j - 1
				break
			}
		}
	}

	return IMAPCommand{Command: "EXPUNGE"}
}

// createFolder generates a CREATE command with
// a randomly generated foldername. The newly created
// folder is appended to the set of folders.
func createFolder(folders *[]Folder) IMAPCommand {

	var arguments []string

	initFoldername := utils.GenerateString(8)

	// Rerandom in case the generated folder name
	// already exists in this session.
	for j := 0; j < len(*folders); j++ {
		if initFoldername == (*folders)[j].Foldername {
			initFoldername = utils.GenerateString(8)
			j = -1
		}
	}

	var messages []Message

	initFolder := Folder{Foldername: initFoldername, Messages: messages}

	*folders = append(*folders, initFolder)
	arguments = append(arguments, initFoldername)
	return IMAPCommand{Command: "CREATE", Arguments: arguments}
}

// deleteFolder generates a DELETE command by deleting
// a random folder from the set of folders. Moreover
// the index of the selected folder is adjusted accordingly.
func deleteFolder(folders *[]Folder, selected *int) IMAPCommand {

	var arguments []string

	folderIndex := rand.Intn(len(*folders))

	for folderIndex == *selected {
		folderIndex = rand.Intn(len(*folders))
	}

	foldername := (*folders)[folderIndex].Foldername

	*folders = append((*folders)[:folderIndex], (*folders)[folderIndex+1:]...)

	if folderIndex < *selected {
		*selected = *selected - 1
	}

	arguments = append(arguments, foldername)
	return IMAPCommand{Command: "DELETE", Arguments: arguments}
}

// selectFolder generates a SELECT command by choosing a random
// folder from the set of folders. Moreover the index of the selected
// folder is adjusted accordingly.
func selectFolder(folders *[]Folder, selected *int) IMAPCommand {

	var arguments []string

	folderIndex := rand.Intn(len(*folders))
	foldername := (*folders)[folderIndex].Foldername

	arguments = append(arguments, foldername)

	*selected = folderIndex
	return IMAPCommand{Command: "SELECT", Arguments: arguments}
}

// appendMsg generates an APPEND command by choosing a random folder
// from the set of folders. A randomly generated message is appended
// to that folder.
func appendMsg(folders *[]Folder) IMAPCommand {

	var arguments []string

	// Choose the folder.
	folderIndex := rand.Intn(len(*folders))

	// Lookup folder name and add it to the arguments list.
	foldername := (*folders)[folderIndex].Foldername
	arguments = append(arguments, foldername)

	// Generate flags of the message.
	flagstring, flags := utils.GenerateFlags()
	arguments = append(arguments, flagstring)

	// Generate date/time string - OPTIONAL.
	arguments = append(arguments, "{310}")

	// Generate message.
	// TODO: Replace with a proper message generator.
	var msg string
	msg = fmt.Sprintf("Date: Mon, 7 Feb 1994 21:52:25 -0800 (PST)\r\nFrom: Fred Foobar <foobar@Blurdybloop.COM>\r\nSubject: afternoon meeting\r\nTo: mooch@owatagu.siam.edu\r\nMessage-Id: <B27397-0100000@Blurdybloop.COM>\r\nMIME-Version: 1.0\r\nContent-Type: TEXT/PLAIN; CHARSET=US-ASCII\r\n\r\nHello Joe, do you think we can meet at 3:30 tomorrow?\r\n")

	arguments = append(arguments, msg)

	(*folders)[folderIndex].Messages = append((*folders)[folderIndex].Messages, Message{Flags: flags})

	return IMAPCommand{Command: "APPEND", Arguments: arguments}

}

// storeMsg generates a STORE command by choosing a random
// message and a random set of flags. The flags of the message
// will be overridden.
func storeMsg(folder *Folder) IMAPCommand {

	var arguments []string

	// Select message.
	msgIndex := rand.Intn(len(folder.Messages))
	arguments = append(arguments, strconv.Itoa(msgIndex+1))

	flagstring, flags := utils.GenerateFlags()
	arguments = append(arguments, flagstring)

	folder.Messages[msgIndex].Flags = flags

	return IMAPCommand{Command: "STORE", Arguments: arguments}
}

// GenerateSession generates a random sequence of IMAPCommands.
// The length of the sequence is between minLength and maxLength.
func GenerateSession(minLength int, maxLength int) []IMAPCommand {

	selected := -1

	var commands []IMAPCommand
	var folders []Folder

	// Define session length.
	sessionLength := rand.Intn((maxLength - minLength)) + minLength

	// Generate the session content.
	for i := 0; i < sessionLength; i++ {

		r := rand.Float64()

		// The following lines represent the allowed IMAP states
		// in a session. Based on the current state of the mailbox,
		// certain IMAP commands might not be allowed e.g.
		// a DELETE is only allowed when there are any folders
		// to be deleted. Hence the following state tree.

		if len(folders) == 0 {

			// We begin with the case where the mailbox is empty.
			// Hence CREATE is the only allowed command.

			commands = append(commands, createFolder(&folders))
		} else {

			// If there are folders in the mailbox, we need
			// to check wheter a folder has been selected.
			// Depending on the selected state, other commands
			// might be allowed.

			if selected == -1 {

				// If the mailbox contains at least one folder and
				// no folder has been selected by SELECT, we allow
				// the following commands:
				// CREATE, DELETE, APPEND, SELECT

				switch {
				case 0.0 <= r && r < 0.25:
					commands = append(commands, createFolder(&folders))
				case 0.25 <= r && r < 0.5:
					commands = append(commands, deleteFolder(&folders, &selected))
				case 0.5 <= r && r < 0.75:
					commands = append(commands, appendMsg(&folders))
				case 0.75 <= r && r < 1.0:
					commands = append(commands, selectFolder(&folders, &selected))
				}

			} else {

				// In this case the mailbox contains at least one folder
				// and a one of these folders has been selected.
				// Next, we need to check wheter there are other folders
				// in the mailbox in order to allow/disallow commands like:
				// DELETE or SELECT.

				if len(folders) == 1 {

					// In case the mailbox contains only one folder and
					// this folder is selected, we need to check wheter
					// there are any messages in the folder in order
					// to allow/disallow the STORE command.

					if len(folders[selected].Messages) == 0 {

						// If there are no messages present in the selected
						// folder, we only allow the following commands:
						// CREATE, APPEND, EXPUNGE

						switch {
						case 0.0 <= r && r < 0.3:
							commands = append(commands, createFolder(&folders))
						case 0.3 <= r && r < 0.9:
							commands = append(commands, appendMsg(&folders))
						case 0.9 <= r && r < 1.0:
							commands = append(commands, expungeFolder(&folders[selected]))
						}

					} else {

						// If there are messages in the selected folder,
						// we can allow STORE as well. Hence the following
						// commands are allowed in this case:
						// CREATE, APPEND, STORE, EXPUNGE
						switch {
						case 0.0 <= r && r < 0.25:
							commands = append(commands, createFolder(&folders))
						case 0.25 <= r && r < 0.5:
							commands = append(commands, appendMsg(&folders))
						case 0.5 <= r && r < 0.75:
							commands = append(commands, storeMsg(&folders[selected]))
						case 0.75 <= r && r < 1.0:
							commands = append(commands, expungeFolder(&folders[selected]))
						}
					}

				} else {

					// In this case the mailbox contains more that one
					// folders and one of these folders is selected.
					// This represents to case with the most variety of
					// IMAP commands. Nevertheless we need to check
					// wheter there are messages in the selected folder
					// in order to allow/disallow the STORE command.

					if len(folders[selected].Messages) == 0 {

						// If there are no messages present, we allow
						// everything except the STORE command:
						// CREATE, DELETE, APPEND, SELECT, EXPUNGE

						switch {
						case 0.0 <= r && r < 0.15:
							commands = append(commands, createFolder(&folders))
						case 0.15 <= r && r < 0.3:
							commands = append(commands, deleteFolder(&folders, &selected))
						case 0.3 <= r && r < 0.6:
							commands = append(commands, appendMsg(&folders))
						case 0.6 <= r && r < 0.9:
							commands = append(commands, selectFolder(&folders, &selected))
						case 0.9 <= r && r < 1.0:
							commands = append(commands, expungeFolder(&folders[selected]))
						}

					} else {

						// In this case we basically allow every IMAP command:
						// CREATE, DELETE, APPEND, STORE, SELECT, EXPUNGE

						switch {
						case 0.0 <= r && r < 0.15:
							commands = append(commands, createFolder(&folders))
						case 0.15 <= r && r < 0.3:
							commands = append(commands, deleteFolder(&folders, &selected))
						case 0.3 <= r && r < 0.5:
							commands = append(commands, appendMsg(&folders))
						case 0.5 <= r && r < 0.75:
							commands = append(commands, storeMsg(&folders[selected]))
						case 0.75 <= r && r < 0.9:
							commands = append(commands, selectFolder(&folders, &selected))
						case 0.9 <= r && r < 1.0:
							commands = append(commands, expungeFolder(&folders[selected]))
						}
					}
				}
			}
		}
	}

	// select the inbox
	var arguments []string
	arguments = append(arguments, "INBOX")
	commands = append(commands, IMAPCommand{Command: "SELECT", Arguments: arguments})

	// Finish session by deleting all created folders.
	for i := 0; i < len(folders); i++ {
		var args []string
		args = append(args, folders[i].Foldername)
		commands = append(commands, IMAPCommand{Command: "DELETE", Arguments: args})
	}

	return commands
}
