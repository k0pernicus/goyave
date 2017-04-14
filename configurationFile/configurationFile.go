/*Package configurationFile represents Encodable/Decodable Golang structures from/to TOML structures.
 *
 *The global structure is ConfigurationFile, which is the simple way to store accurtely your local informations.
 *
 */
package configurationFile

import (
	"bytes"

	"errors"

	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/k0pernicus/goyave/consts"
)

/*ConfigurationFile represents the TOML structure of the Goyave configuration file.
 *
 *The structure of the configuration file is:
 *	Author:
 *		The name of the user.
 *	VisibleRepositories:
 *		A list of local visible git repositories.
 *	HiddenRepositories:
 *		A list of local ignored git repositories.
 *	Groups:
 *		A list of groups.
 */
type ConfigurationFile struct {
	Author              string
	Local               LocalInformations `toml:"local"`
	VisibleRepositories []GitRepository   `toml:"visible"`
	HiddenRepositories  []GitRepository   `toml:"hidden"`
	Groups              []Group           `toml:"group"`
}

/*AddVisibleRepository adds a given git repo path as a visible repository
 *
 *If the repository already exists in the VisibleRepository field, the method throws an error: RepositoryAlreadyExists.
 *Else, the repository is append to the VisibleRepository field, and the method returns nil.
 */
func (c *ConfigurationFile) AddVisibleRepository(path string) error {
	for _, registeredRepository := range c.VisibleRepositories {
		if registeredRepository.Path == path {
			return errors.New(consts.RepositoryAlreadyExists)
		}
	}
	var newVisibleRepository = GitRepository{
		Name: filepath.Ext(path),
		Path: path,
	}
	c.VisibleRepositories = append(c.VisibleRepositories, newVisibleRepository)
	return nil
}

/*AddHiddenRepository adds a given git repo path as an hidden repository
 *
 *If the repository already exists in the HiddenRepository field, the method throws an error: RepositoryAlreadyExists.
 *Else, the repository is append to the HiddenRepository field, and the method returns nil.
 */
func (c *ConfigurationFile) AddHiddenRepository(path string) error {
	for _, registeredRepository := range c.HiddenRepositories {
		if registeredRepository.Path == path {
			return errors.New(consts.RepositoryAlreadyExists)
		}
	}
	var newHiddenRepository = GitRepository{
		Name: filepath.Ext(path),
		Path: path,
	}
	c.HiddenRepositories = append(c.HiddenRepositories, newHiddenRepository)
	return nil
}

/*GitRepository represents the structure of a local git repository.
 *
 *Properties of this structure are:
 *	Name:
 * 		The custom name of the repository.
 *	Path:
 *		The path of the repository.
 */
type GitRepository struct {
	Name string
	Path string
}

/*Group represents a group of git repositories.
 *
 *The structure of a Group type is:
 *
 *	Name:
 *		The group name.
 *	Repositories:
 *		A list of git repositories id, tagged in the group.
 */
type Group struct {
	Name         string
	Repositories []string
}

/*LocalInformations represents your local configuration of Goyave.
 *
 *The structure contains:
 *	Group:
 *		The current group name.
 */
type LocalInformations struct {
	Group string
}

/*Decode is a function to decode an entire string (which is the content of a given TOML file) to a ConfigurationFile structure.
 */
func Decode(data string, localStructure *ConfigurationFile) error {
	_, err := toml.Decode(data, localStructure)
	return err
}

/*Encode is a function to encode a ConfigurationFile structure to a byffer of bytes.
 */
func Encode(localStructure *ConfigurationFile, buffer *bytes.Buffer) error {
	return toml.NewEncoder(buffer).Encode(localStructure)
}
