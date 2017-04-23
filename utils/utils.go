package utils

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/k0pernicus/goyave/consts"
)

/*IsGitRepository returns if the path, given as an argument, is a git repository or not.
 *This function returns a boolean value: true if the pathdir pointed to a git repository, else false.
 */
func IsGitRepository(pathdir string) bool {
	if filepath.Base(pathdir) != consts.GitFileName {
		pathdir = filepath.Join(pathdir, consts.GitFileName)
	}
	file, err := os.Open(pathdir)
	if err != nil {
		return false
	}
	_, err = file.Stat()
	return !os.IsNotExist(err)
}

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
