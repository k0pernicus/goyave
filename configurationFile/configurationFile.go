/*
Package configurationFile represents Encodable/Decodable Golang structures from/to TOML structures.

The global structure is ConfigurationFile, which is the simple way to store accurtely your local informations.

*/
package configurationFile

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
)

/*
ConfigurationFile represents the TOML structure of the Goyave configuration file.

The structure of the configuration file is:

	Author:
		The name of the user.
	VisibleRepositories:
		A list of local visible git repositories.
	HiddenRepositories:
		A list of local ignored git repositories.
	Groups:
		A list of groups.
*/
type ConfigurationFile struct {
	Author              string
	Local               LocalInformations `toml:"local"`
	VisibleRepositories []GitRepository   `toml:"visible"`
	HiddenRepositories  []GitRepository   `toml:"hidden"`
	Groups              []Group           `toml:"group"`
}

/*
GitRepository represents the structure of a local git repository.

Properties of this structure are:

	Name:
		The custom name of the repository.
	Path:
		The path of the repository.
*/
type GitRepository struct {
	Name string
	Path string
}

/*
Group represents a group of git repositories.

The structure of a Group type is:

	Name:
		The group name.
	Repositories:
		A list of git repositories id, tagged in the group.
*/
type Group struct {
	Name         string
	Repositories []string
}

/*
LocalInformations represents your local configuration of Goyave.

The structure contains:

	Group:
		The current group name.
*/
type LocalInformations struct {
	Group string
}

/*
Decode is a function to decode an entire string (which is the content of a given TOML file) to a ConfigurationFile structure.
*/
func Decode(data string, localStructure *ConfigurationFile) error {
	_, err := toml.Decode(data, localStructure)
	return err
}

/*
Encode is a function to encode a ConfigurationFile structure to a byffer of bytes.
*/
func Encode(localStructure *ConfigurationFile, buffer *bytes.Buffer) error {
	return toml.NewEncoder(buffer).Encode(localStructure)
}

/*
Open returns a pointer for the configuration file.

By default, the mod of this configuration file is 0755.
*/
func Open(configurationFilePath *string) (*os.File, error) {
	return os.OpenFile(*configurationFilePath, os.O_RDWR|os.O_CREATE, 0755)
}
