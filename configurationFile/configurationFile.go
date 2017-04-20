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

	"fmt"

	"os/user"

	"github.com/BurntSushi/toml"
	"github.com/k0pernicus/goyave/consts"
	"github.com/k0pernicus/goyave/utils"
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

/*Default returns a default ConfigurationFile structure.
 */
func Default() *ConfigurationFile {
	usr, _ := user.Current()
	return &ConfigurationFile{
		Author: usr.Username,
		Local: LocalInformations{
			DefaultEntry: consts.VisibleFlag,
			Group:        utils.GetLocalhost(),
		},
	}
}

/*GetDefaultEntry returns the default entry to store a new git repository.
 *This methods returns HiddenFlag, or VisibleFlag
 */
func (c *ConfigurationFile) GetDefaultEntry() (string, error) {
	defaultEntry := c.Local.DefaultEntry
	if defaultEntry != consts.VisibleFlag && defaultEntry != consts.HiddenFlag {
		return consts.VisibleFlag, fmt.Errorf("the default entry is not set to %s or %s. Please check your configuration file", consts.HiddenFlag, consts.VisibleFlag)
	}
	return defaultEntry, nil
}

/*AddRepository will add a single path in the TOML's target if the path does not exists.
 *
 *This method uses both methods addVisibleRepository and addHiddenRepository.
 */
func (c *ConfigurationFile) AddRepository(path string, target string) error {
	if target == consts.VisibleFlag {
		return c.addVisibleRepository(path)
	}
	if target == consts.HiddenFlag {
		return c.addHiddenRepository(path)
	}
	return errors.New("The target does not exists.")
}

/*addVisibleRepository adds a given git repo path as a visible repository
 *
 *If the repository already exists in the VisibleRepository field, the method throws an error: RepositoryAlreadyExists.
 *Else, the repository is append to the VisibleRepository field, and the method returns nil.
 */
func (c *ConfigurationFile) addVisibleRepository(path string) error {
	for _, registeredRepository := range c.VisibleRepositories {
		if registeredRepository.Path == path {
			return errors.New(consts.RepositoryAlreadyExists)
		}
	}
	var newVisibleRepository = GitRepository{
		Name: filepath.Base(path),
		Path: path,
	}
	c.VisibleRepositories = append(c.VisibleRepositories, newVisibleRepository)
	return nil
}

/*addHiddenRepository adds a given git repo path as an hidden repository
 *
 *If the repository already exists in the HiddenRepository field, the method throws an error: RepositoryAlreadyExists.
 *Else, the repository is append to the HiddenRepository field, and the method returns nil.
 */
func (c *ConfigurationFile) addHiddenRepository(path string) error {
	for _, registeredRepository := range c.HiddenRepositories {
		if registeredRepository.Path == path {
			return errors.New(consts.RepositoryAlreadyExists)
		}
	}
	var newHiddenRepository = GitRepository{
		Name: filepath.Base(path),
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
 *  DefaultEntry:
 *		The default entry to store a git repository (hidden or visible).
 *	Group:
 *		The current group name.
 */
type LocalInformations struct {
	DefaultEntry string
	Group        string
}

/*DecodeString is a function to decode an entire string (which is the content of a given TOML file) to a ConfigurationFile structure.
 */
func (c *ConfigurationFile) DecodeString(data string) error {
	_, err := toml.Decode(data, c)
	return err
}

/*DecodeBytesArray is a function to decode an entire string (which is the content of a given TOML file) to a ConfigurationFile structure.
 */
func (c *ConfigurationFile) DecodeBytesArray(data []byte) error {
	_, err := toml.Decode(string(data[:]), c)
	return err
}

/*Encode is a function to encode a ConfigurationFile structure to a byffer of bytes.
 */
func (c *ConfigurationFile) Encode(buffer *bytes.Buffer) error {
	return toml.NewEncoder(buffer).Encode(c)
}
