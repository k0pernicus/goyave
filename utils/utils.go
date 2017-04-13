package utils

import (
	"os"
	"os/user"
)

/*
GetUserHomeDir returns the home directory of the current user.
*/
func GetUserHomeDir() string {
	usr, err := user.Current()
	// If the current user cannot be reached, get the HOME environment variable
	if err != nil {
		return os.Getenv("$HOME")
	}
	return usr.HomeDir
}
