package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"sync"

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

	/*stateCmd is a subcommand to list the state of each local git repository.
	 */
	var stateCmd = &cobra.Command{
		Use:   "state",
		Short: "Get the state of each local git repository",
		Run: func(cmd *cobra.Command, args []string) {
			for _, gitStruct := range configurationFileStructure.VisibleRepositories {
				gitStruct.Init()
				gitStruct.GitObject.Status()
			}
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

	rootCmd.AddCommand(addCmd, crawlCmd, stateCmd, switchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
