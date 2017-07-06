package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"sync"

	"sort"

	"github.com/BurntSushi/toml"
	"github.com/k0pernicus/goyave/configurationFile"
	"github.com/k0pernicus/goyave/consts"
	"github.com/k0pernicus/goyave/traces"
	"github.com/k0pernicus/goyave/utils"
	"github.com/k0pernicus/goyave/walk"
	"github.com/spf13/cobra"
)

var configurationFileStructure configurationFile.ConfigurationFile
var configurationFilePath string
var userHomeDir string

/*initialize get the configuration file existing in the system (or create it), and return
 *a pointer to his content.
 */
func initialize(configurationFileStructure *configurationFile.ConfigurationFile) {
	// Initialize all different traces structures
	traces.InitTraces(os.Stdout, os.Stderr, os.Stdout, os.Stdout)

	// Get the user home directory
	userHomeDir = utils.GetUserHomeDir()
	if len(userHomeDir) == 0 {
		log.Fatalf("Cannot get the user home dir.\n")
	}
	// Set the configuration path file
	configurationFilePath = path.Join(userHomeDir, consts.ConfigurationFileName)
	filePointer, err := os.OpenFile(configurationFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Cannot open the file %s, due to error '%s'.\n", configurationFilePath, err)
	}
	defer filePointer.Close()
	var bytesArray []byte
	// Get the content of the goyave configuration file
	configurationFile.GetConfigurationFileContent(filePointer, &bytesArray)
	if _, err = toml.Decode(string(bytesArray[:]), configurationFileStructure); err != nil {
		log.Fatalln(err)
	}
}

/*kill saves the current state of the configuration structure in the configuration file
 */
func kill() {
	var outputBuffer bytes.Buffer
	if err := configurationFileStructure.Encode(&outputBuffer); err != nil {
		log.Fatalln("Cannot save the current configurationFile structure!")
	}
	if err := ioutil.WriteFile(configurationFilePath, outputBuffer.Bytes(), 0777); err != nil {
		log.Fatalln("Cannot access to your file to save the configurationFile structure!")
	}
}

func main() {

	/*rootCmd defines the global app, and some actions to run before and after the command running
	 */
	var rootCmd = &cobra.Command{
		Use:   "goyave",
		Short: "Goyave is a tool to take a look at your local git repositories",
		// Initialize the structure
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initialize(&configurationFileStructure)
		},
		// Save the current configuration file structure, in the configuration file
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			kill()
		},
	}

	/*addCmd is a subcommand to add the current working directory as a VISIBLE one
	 */
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add the current path as a VISIBLE repository",
		Run: func(cmd *cobra.Command, args []string) {
			// Get the path where the command has been executed
			currentDir, err := os.Getwd()
			if err != nil {
				log.Fatalln("There was a problem retrieving the current directory")
			}
			if !utils.IsGitRepository(currentDir) {
				log.Fatalf("%s is not a git repository!\n", currentDir)
			}
			// If the path is/contains a .git directory, add this one as a VISIBLE repository
			if err := configurationFileStructure.AddRepository(currentDir, consts.VisibleFlag); err != nil {
				traces.WarningTracer.Printf("[%s] %s\n", currentDir, err)
			}
		},
	}

	/*crawlCmd is a subcommand to crawl your hard drive in order to get and save new git repositories
	 */
	var crawlCmd = &cobra.Command{
		Use:   "crawl",
		Short: "Crawl the hard drive in order to find git repositories",
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			// Get all git paths, and display them
			gitPaths, err := walk.RetrieveGitRepositories(userHomeDir)
			if err != nil {
				log.Fatalf("There was an error retrieving your git repositories: '%s'\n", err)
			}
			// For each git repository, check if it exists, and if not add it to the default target visibility
			for _, gitPath := range gitPaths {
				wg.Add(1)
				go func(gitPath string) {
					defer wg.Done()
					if err := configurationFileStructure.AddRepository(gitPath, configurationFileStructure.Local.DefaultTarget); err != nil {
						traces.WarningTracer.Printf("[%s] %s\n", gitPath, err)
					}
				}(gitPath)
			}
			wg.Wait()
		},
	}

	/*pathCmd is a subcommand to get the path of a given git repository.
	 *This subcommand is useful to change directory, like `cd $(goyave path mygitrepo)`
	 */
	var pathCmd = &cobra.Command{
		Use:   "path",
		Short: "Get the path of a given repository, if this one exists",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalln("Needs a repository name!")
			}
			repo := args[0]
			repoPath := configurationFileStructure.GetPathFromRepository(repo)
			if repoPath != "" {
				fmt.Println(repoPath)
			} else {
				log.Fatalf("The repository %s does not exists!\n", repo)
			}
		},
	}

	/*stateCmd is a subcommand to list the state of each local git repository.
	 */
	var stateCmd = &cobra.Command{
		Use:     "state",
		Example: "goyave state\ngoyave state myRepositoryName\ngoyave state myRepositoryName1 myRepositoryName2",
		Short:   "Get the state of each local visible git repository",
		Long:    "Check only visible git repositories.\nIf some repository names have been setted, goyave will only check those repositories, otherwise it checks all visible repositories of your system.",
		Run: func(cmd *cobra.Command, args []string) {
			var gitStructs []configurationFile.GitRepository
			if len(args) == 0 {
				gitStructs = configurationFileStructure.VisibleRepositories
			} else {
				// Sort visible repositories by name
				sort.Sort(configurationFile.ByName(configurationFileStructure.VisibleRepositories))
				repositoriesListLength := len(configurationFileStructure.VisibleRepositories)
				// Looking for given repository names - if the looking one does not exists, let the function prints a warning message.
				for _, repositoryName := range args {
					repositoryIndex := sort.Search(repositoriesListLength, func(i int) bool { return configurationFileStructure.VisibleRepositories[i].Name >= repositoryName })
					if repositoryIndex != repositoriesListLength {
						gitStructs = append(gitStructs, configurationFileStructure.VisibleRepositories[repositoryIndex])
					} else {
						traces.WarningTracer.Printf("%s cannot be found in your visible repositories!\n", repositoryName)
					}
				}
			}
			var wg sync.WaitGroup
			for _, gitStruct := range gitStructs {
				wg.Add(1)
				go func(gitStruct configurationFile.GitRepository) {
					defer wg.Done()
					gitStruct.Init()
					gitStruct.GitObject.Status()
				}(gitStruct)
			}
			wg.Wait()
		},
	}

	/*switchCmd is a subcommand to switch the visibility of the current git repository.
	 */
	var switchCmd = &cobra.Command{
		Use:   "switch",
		Short: "Switch the visibility of the current git repository (given by the current path)",
		Run: func(cmd *cobra.Command, args []string) {
			// Get the path where the command has been executed
			currentDir, err := os.Getwd()
			if err != nil {
				log.Fatalln("There was a problem retrieving the current directory")
			}
			if err := configurationFileStructure.RemoveRepositoryFromSlice(currentDir, consts.VisibleFlag); err == nil {
				configurationFileStructure.AddRepository(currentDir, consts.HiddenFlag)
				traces.InfoTracer.Printf("%s has been set to an hidden repository!", currentDir)
				return
			}
			if err := configurationFileStructure.RemoveRepositoryFromSlice(currentDir, consts.HiddenFlag); err == nil {
				configurationFileStructure.AddRepository(currentDir, consts.VisibleFlag)
				traces.InfoTracer.Printf("%s has been set to a visible repository!", currentDir)
				return
			}
			log.Fatalf("The repository %s is not saved as a VISIBLE or HIDDEN repository! Please to add it before.\n", currentDir)
		},
	}

	rootCmd.AddCommand(addCmd, crawlCmd, pathCmd, stateCmd, switchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
