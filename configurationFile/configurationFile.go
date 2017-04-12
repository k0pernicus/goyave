/*
Package configurationFile represents Encodable/Decodable Golang structures from/to TOML structures.

The global structure is ConfigurationFile, which is the simple way to store accurtely your local informations.

*/
package configurationFile

import "github.com/BurntSushi/toml"
import "log"
import "bytes"

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
	Name         string `toml:"name"`
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

If the TOML content cannot be decoded, this function throw an error.
*/
func Decode(data string, localStructure *ConfigurationFile) {
	if _, err := toml.Decode(data, localStructure); err != nil {
		log.Fatal(err)
	}
}

/*
Encode is a function to encode a ConfigurationFile structure to a byffer of bytes.

This function can throw an error if the TOML encoder failed.
*/
func Encode(localStructure *ConfigurationFile, buffer *bytes.Buffer) {
	if err := toml.NewEncoder(buffer).Encode(localStructure); err != nil {
		log.Fatal(err)
	}
}
