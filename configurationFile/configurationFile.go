/*Package configurationFile represents Encodable/Decodable Golang structures from/to TOML structures.
 *
 *The global structure is ConfigurationFile, which is the simple way to store accurtely your local informations.
 *
 */
package configurationFile

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"

	"fmt"

	"path/filepath"

	"sync"

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
		cLocalhost := utils.GetLocalhost()
		cUser, err := user.Current()
		var cUserName string
		if err != nil {
			cUserName = consts.DefaultUserName
		} else {
			cUserName = cUser.Username
		}
		defaultStructure := Default(cUserName, cLocalhost)
		defaultStructure.Encode(&fileBuffer)
		*bytesArray = fileBuffer.Bytes()
	} else {
		b, _ := ioutil.ReadAll(filePointer)
		*bytesArray = b
	}
}

/*ConfigurationFile represents the TOML structure of the Goyave configuration file
 *
 *Properties:
 *	Author:
 *		The name of the user
 *  Local:
 *		Local informations
 *  Repositories:
 *		Local git repositories
 *	VisibleRepositories:
 *		A list of local ** visible ** git repositories (** used localy **)
 *	Groups:
 *		A list of groups
 *	locker:
 *		Mutex to perform concurrent RW on map data structures
 */
type ConfigurationFile struct {
	Author              string
	Local               LocalInformations        `toml:"local"`
	Repositories        map[string]GitRepository `toml:"repositories"`
	VisibleRepositories VisibleRepositories      `toml:"-"`
	Groups              map[string]Group         `toml:"group"`
	locker              sync.RWMutex             `toml:"-"`
}

/*Default is a constructor for ConfigurationFile
 *
 *Parameters:
 *	author:
 *		The name of the user
 *  hostname:
 *		The machine hostname
 */
func Default(author string, hostname string) *ConfigurationFile {
	return &ConfigurationFile{
		Author: author,
		Local: LocalInformations{
			DefaultTarget: consts.VisibleFlag,
			Group:         hostname,
		},
		Groups: map[string]Group{
			hostname: []string{},
		},
	}
}

/*AddRepository append the given repository to the list of local repositories, if it does not exists
 */
func (c *ConfigurationFile) AddRepository(path, target string) error {
	name := filepath.Base(path)
	hostname := utils.GetLocalhost()
	c.locker.Lock()
	defer c.locker.Unlock()
	robj, ok := c.Repositories[name]
	// If the repository exists and the path is ok, stop
	if ok && robj.Paths[hostname].Path == path {
		return nil
	}
	// Initialize the new GroupPath structure
	cgroup := GroupPath{
		Name: name,
		Path: path,
	}
	// If the repository exists but the path is not ok, update it
	if ok {
		robj.Paths[hostname] = cgroup
		return nil
	}
	// Otherwise, create a new GitRepository structure, and append it in the Repositories field
	c.Repositories[name] = GitRepository{
		Name: name,
		Paths: map[string]GroupPath{
			hostname: cgroup,
		},
		URL: gitManip.GetRemoteURL(path),
	}
	// If the user wants to add automatically new repositories as repositories to "follow", change
	// his flag as a "visible" repository
	if target == consts.VisibleFlag {
		c.Groups[hostname] = append(c.Groups[hostname], name)
	}
	return nil
}

/*GetPath returns the local path file, for a given repository
 */
func (c *ConfigurationFile) GetPath(repository string) (string, bool) {
	gobj, ok := c.VisibleRepositories[repository]
	return gobj, ok
}

/*Process initializes useful fields in the data structure
 */
func (c *ConfigurationFile) Process() error {
	// If the configuration file is new, initialize the map and finish here
	if c.Repositories == nil {
		c.Repositories = make(map[string]GitRepository)
		return nil
	}
	// Otherwise, initialize useful fields
	hostname := utils.GetLocalhost()
	vrepositories, ok := c.Groups[hostname]
	if !ok {
		return fmt.Errorf("the hostname %s has not been found - please to launch 'crawl' before", hostname)
	}
	c.VisibleRepositories = make(VisibleRepositories)
	for _, repository := range vrepositories {
		c.VisibleRepositories[repository] = c.Repositories[repository].Paths[hostname].Path
	}
	return nil
}

/*VisibleRepositories is a map structure to store, for each repository name (and the hostname), the associated path
 */
type VisibleRepositories map[string]string

/*Method that returns if a repository, identified by his name and his path (optional), exists in the given structure
 * If path is empty (empty string), the function will only check the name
 */
func (v VisibleRepositories) exists(name, path string) bool {
	_, ok := v[name]
	if !ok || path == "" {
		return ok
	}
	return v[name] == path
}

/*GitRepository represents the structure of a local git repository
 *
 *Properties:
 *	Name:
 * 		The custom name of the repository
 *  Paths:
 *		Path per group name
 *	URL:
 *		The remote URL of the repository (from origin)
 */
type GitRepository struct {
	Name  string               `toml:"name"`
	Paths map[string]GroupPath `toml:"paths"`
	URL   string               `toml:"url"`
}

/*GroupPath represents the structure of a local path, using a given group
 *
 *Properties:
 *  Name:
 *		The name of the local git repository
 *	Path:
 *		A string that points to the local git repository
 */
type GroupPath struct {
	Name string
	Path string
}

/*Group represents a group of git repositories names
 */
type Group []string

/*LocalInformations represents your local configuration of Goyave
 *
 *Properties:
 *  DefaultEntry:
 *		The default entry to store a git repository (hidden or visible)
 *	Group:
 *		The current group name.
 */
type LocalInformations struct {
	DefaultTarget string
	Group         string
}

/*DecodeString is a function to decode an entire string (which is the content of a given TOML file) to a ConfigurationFile structure
 */
func DecodeString(c *ConfigurationFile, data string) error {
	_, err := toml.Decode(data, *c)
	return err
}

/*DecodeBytesArray is a function to decode an entire string (which is the content of a given TOML file) to a ConfigurationFile structure
 */
func DecodeBytesArray(c *ConfigurationFile, data []byte) error {
	_, err := toml.Decode(string(data[:]), *c)
	return err
}

/*Encode is a function to encode a ConfigurationFile structure to a byffer of bytes
 */
func (c *ConfigurationFile) Encode(buffer *bytes.Buffer) error {
	return toml.NewEncoder(buffer).Encode(c)
}
