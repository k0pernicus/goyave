/*Package configurationFile represents Encodable/Decodable Golang structures from/to TOML structures.
 *
 *The global structure is ConfigurationFile, which is the simple way to store accurtely your local informations.
 *
 */
package configurationFile

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"

	"fmt"

	"os/user"

	"os"

	"strings"

	"github.com/BurntSushi/toml"
	"github.com/k0pernicus/goyave/consts"
	"github.com/k0pernicus/goyave/gitManip"
	"github.com/k0pernicus/goyave/traces"
	"github.com/k0pernicus/goyave/utils"
)

/*GetConfigurationFileContent get the content of the local configuration file.
 *If no configuration file has been found, create a default one and set the bytes array.
 */
func GetConfigurationFileContent(filePointer *os.File, bytesArray *[]byte) {
	fileState, err := filePointer.Stat()
	// If the file is empty, get the default structure and save it
	if err != nil || fileState.Size() == 0 {
		traces.WarningTracer.Println("No (or empty) configuration file - creating default one...")
		var fileBuffer bytes.Buffer
		defaultStructure := Default()
		defaultStructure.Encode(&fileBuffer)
		*bytesArray = fileBuffer.Bytes()
	} else {
		b, _ := ioutil.ReadAll(filePointer)
		*bytesArray = b
	}
}

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
			DefaultTarget: consts.VisibleFlag,
			Group:         utils.GetLocalhost(),
		},
	}
}

/*GetDefaultEntry returns the default entry to store a new git repository.
 *This methods returns HiddenFlag, or VisibleFlag
 */
func (c *ConfigurationFile) GetDefaultEntry() (string, error) {
	defaultTarget := c.Local.DefaultTarget
	if defaultTarget != consts.VisibleFlag && defaultTarget != consts.HiddenFlag {
		return consts.VisibleFlag, fmt.Errorf("the default target is not set to %s or %s. Please check your configuration file", consts.HiddenFlag, consts.VisibleFlag)
	}
	return defaultTarget, nil
}

/*AddRepository will add a single path in the TOML's target if the path does not exists.
 *This method uses both methods addVisibleRepository and addHiddenRepository.
 */
func (c *ConfigurationFile) AddRepository(path string, target string) error {
	if target == consts.VisibleFlag {
		return c.addVisibleRepository(path)
	}
	if target == consts.HiddenFlag {
		return c.addHiddenRepository(path)
	}
	return errors.New("the target does not exists")
}

/*addVisibleRepository adds a given git repo path as a visible repository.
 *If the repository already exists in the VisibleRepository field, the method throws an error: RepositoryAlreadyExists.
 *Else, the repository is append to the VisibleRepository field, and the method returns nil.
 */
func (c *ConfigurationFile) addVisibleRepository(path string) error {
	for _, registeredRepository := range c.VisibleRepositories {
		if registeredRepository.Path == path {
			return errors.New(consts.RepositoryAlreadyExists)
		}
	}
	var newVisibleRepository = NewGitRepository(filepath.Base(path), path)
	c.VisibleRepositories = append(c.VisibleRepositories, newVisibleRepository)
	return nil
}

/*addHiddenRepository adds a given git repo path as an hidden repository.
 *If the repository already exists in the HiddenRepository field, the method throws an error: RepositoryAlreadyExists.
 *Else, the repository is append to the HiddenRepository field, and the method returns nil.
 */
func (c *ConfigurationFile) addHiddenRepository(path string) error {
	for _, registeredRepository := range c.HiddenRepositories {
		if registeredRepository.Path == path {
			return errors.New(consts.RepositoryAlreadyExists)
		}
	}
	var newHiddenRepository = NewGitRepository(filepath.Base(path), path)
	c.HiddenRepositories = append(c.HiddenRepositories, newHiddenRepository)
	return nil
}

/*GetPathFromRepository returns the path of a given git repository name.
 *If the repository does not exists, it returns an empty string.
 */
func (c *ConfigurationFile) GetPathFromRepository(target string) string {
	for _, registeredRepository := range c.VisibleRepositories {
		if registeredRepository.Name == target {
			return registeredRepository.Path
		}
	}
	return ""
}

/*RemoveRepositoryFromSlice returns a new slice without the corresponding element (here, a string).
 *If the element is not found, this method returns an error.
 */
func (c *ConfigurationFile) RemoveRepositoryFromSlice(path string, slice string) error {
	// The code below is the same for visible and hidden repositories - need to refactor the code later
	if slice == consts.VisibleFlag {
		sliceIndex := utils.SliceIndex(len(c.VisibleRepositories), func(i int) bool { return c.VisibleRepositories[i].Path == path })
		if sliceIndex != -1 {
			c.VisibleRepositories = append(c.VisibleRepositories[:sliceIndex], c.VisibleRepositories[sliceIndex+1:]...)
			return nil
		}
		return errors.New(consts.ItemIsNotInSlice)
	}
	sliceIndex := utils.SliceIndex(len(c.HiddenRepositories), func(i int) bool { return c.HiddenRepositories[i].Path == path })
	if sliceIndex != -1 {
		c.HiddenRepositories = append(c.HiddenRepositories[:sliceIndex], c.HiddenRepositories[sliceIndex+1:]...)
		return nil
	}
	return errors.New(consts.ItemIsNotInSlice)
}

/*GitRepository represents the structure of a local git repository.
 *
 *Properties of this structure are:
 *	GitObject:
 *		A reference to a git structure that represents the repository - ignored in the TOML file.
 *	Name:
 * 		The custom name of the repository.
 *	Path:
 *		The path of the repository.
 *	URL:
 *		The remote URL of the repository (from origin).
 */
type GitRepository struct {
	GitObject *gitManip.GitObject `toml:"-"`
	Name      string
	Path      string
	URL       string
}

/*ByName implements sort.Interface for []GitRepository based on the Name field.
 */
type ByName []GitRepository

/*Len returns the length of the ByName type object.
 */
func (g ByName) Len() int { return len(g) }

/*Swap swaps two objects in the same array.
 */
func (g ByName) Swap(i, j int) { g[i], g[j] = g[j], g[i] }

/*Less returns True if the first element is lower than the second one (alphabetic order).
 */
func (g ByName) Less(i, j int) bool { return strings.Compare(g[i].Name, g[j].Name) == -1 }

/*NewGitRepository instantiates the GitRepository struct, based on the path information.
 */
func NewGitRepository(name, path string) GitRepository {
	gitObject := gitManip.New(path)
	return GitRepository{
		GitObject: gitObject,
		Name:      name,
		Path:      path,
		URL:       gitObject.GetRemoteURL(),
	}
}

/*Init (re)initializes the GitObject structure
 */
func (g *GitRepository) Init() {
	g.GitObject = gitManip.New(g.Path)
}

/*isExists check if the current path of the git repository is correct or not,
 *and if the current repository exists again or not.
 *This methods returns a boolean value.
 */
func (g *GitRepository) isExists() bool {
	_, err := os.Stat(g.Path)
	return os.IsNotExist(err)
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
	DefaultTarget string
	Group         string
}

/*DecodeString is a function to decode an entire string (which is the content of a given TOML file) to a ConfigurationFile structure.
 */
func DecodeString(c *ConfigurationFile, data string) error {
	_, err := toml.Decode(data, *c)
	return err
}

/*DecodeBytesArray is a function to decode an entire string (which is the content of a given TOML file) to a ConfigurationFile structure.
 */
func DecodeBytesArray(c *ConfigurationFile, data []byte) error {
	_, err := toml.Decode(string(data[:]), *c)
	return err
}

/*Encode is a function to encode a ConfigurationFile structure to a byffer of bytes.
 */
func (c *ConfigurationFile) Encode(buffer *bytes.Buffer) error {
	return toml.NewEncoder(buffer).Encode(c)
}
