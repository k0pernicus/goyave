package utils

import (
	"os"
	"os/user"
)

/*GetUserHomeDir returns the home directory of the current user.
 */
func GetUserHomeDir() string {
	usr, err := user.Current()
	// If the current user cannot be reached, get the HOME environment variable
	if err != nil {
		return os.Getenv("$HOME")
	}
	return usr.HomeDir
}

/*GetLocalhost returns the localhost name of the current computer.
 *If there is an error, it returns a default string.
 */
func GetLocalhost() string {
	lhost, err := os.Hostname()
	if err != nil {
		return "DefaultHostname"
	}
	return lhost
}
