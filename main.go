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
	"github.com/k0pernicus/goyave/gitManip"
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
		log.Fatalf("cant get the user home dir\n")
	}
	// Set the configuration path file
	configurationFilePath = path.Join(userHomeDir, consts.ConfigurationFileName)
	filePointer, err := os.OpenFile(configurationFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("cant open the file %s, due to error '%s'\n", configurationFilePath, err)
	}
	defer filePointer.Close()
	var bytesArray []byte
	// Get the content of the goyave configuration file
	configurationFile.GetConfigurationFileContent(filePointer, &bytesArray)
	if _, err = toml.Decode(string(bytesArray[:]), configurationFileStructure); err != nil {
		log.Fatalln(err)
	}
	if err := configurationFileStructure.Process(); err != nil {
		log.Fatalln(err)
	}
}

/*kill saves the current state of the configuration structure in the configuration file
 */
func kill() {
	var outputBuffer bytes.Buffer
	if err := configurationFileStructure.Encode(&outputBuffer); err != nil {
		log.Fatalln("can't save the current configurationFile structure")
	}
	if err := ioutil.WriteFile(configurationFilePath, outputBuffer.Bytes(), 0777); err != nil {
		log.Fatalln("can't access to your file to save the configurationFile structure")
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

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Init",
		Run: func(cmd *cobra.Command, args []string) {
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
				log.Fatalf("there was an error retrieving your git repositories: '%s'\n", err)
			}
			// For each git repository, check if it exists, and if not add it to the default target visibility
			for _, gitPath := range gitPaths {
				wg.Add(1)
				go func(gitPath string) {
					defer wg.Done()
					if !utils.IsGitRepository(gitPath) {
						configurationFileStructure.AddRepository(gitPath, configurationFileStructure.Local.DefaultTarget)
					}
				}(gitPath)
			}
			wg.Wait()
		},
	}

	/*loadCmd permits to load visible repositories from the goyave configuration file
	 */
	// var loadCmd = &cobra.Command{
	// 	Use:   "load",
	// 	Short: "Load the configuration file to restore your previous work space",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		currentLocalhost := utils.GetLocalhost()
	// 		fmt.Printf("Current localhost is %s\n", currentLocalhost)
	// 		configGroups := configurationFileStructure.Groups
	// 		var visibleRepositories []string
	// 		for {
	// 			index := utils.SliceIndex(len(configGroups), func(i int) bool { return configGroups[i].Name == currentLocalhost })
	// 			if index == -1 {
	// 				traces.WarningTracer.Printf("Your current local host (%s) has not been found!", currentLocalhost)
	// 				fmt.Println("Please to choose one of those, to load the configuration file:")
	// 				for _, group := range configGroups {
	// 					fmt.Printf("\t%s\n", group.Name)
	// 				}
	// 				scanner := bufio.NewScanner(os.Stdin)
	// 				currentLocalhost = scanner.Text()
	// 				continue
	// 			} else {
	// 				visibleRepositories = configurationFileStructure.Groups[index].VisibleRepositories
	// 			}
	// 			break
	// 		}
	// 		traces.InfoTracer.Printf("Importing configuration from group %s\n", currentLocalhost)
	// 		for _, visibleRepository := range visibleRepositories {
	// 			traces.InfoTracer.Printf("* Importing %s...\n", visibleRepository)
	// 			index := utils.SliceIndex(len(configurationFileStructure.Repositories), func(i int) bool { return configurationFileStructure.Repositories[i].Name == visibleRepository })
	// 			// Check the local path, and the remote URL
	// 			if index == -1 {
	// 				traces.WarningTracer.Printf("\tThe repository \"%s\" does not exists in your configuration file.\n", visibleRepository)
	// 				continue
	// 			}
	// 			// Check if the repository exists locally
	// 			pathRepository, URLRepository := configurationFileStructure.Repositories[index].Path, configurationFileStructure.Repositories[index].URL
	// 			if _, err := os.Stat(pathRepository); err == nil {
	// 				traces.InfoTracer.Printf("\tThe repository \"%s\" already exists as a local git repository.\n", visibleRepository)
	// 				continue
	// 			}
	// 			// If it does not exists, clone it
	// 			if err := gitManip.Clone(pathRepository, URLRepository); err != nil {
	// 				traces.ErrorTracer.Printf("\tThe repository \"%s\" can't be cloned: %s\n", visibleRepository, err)
	// 			} else {
	// 				traces.InfoTracer.Printf("\tThe repository \"%s\" has been successfully cloned!\n", visibleRepository)
	// 			}
	// 		}
	// 	},
	// }

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
			repoPath, found := configurationFileStructure.GetPath(repo)
			if !found {
				log.Fatalf("repository %s not found\n", repo)
			} else {
				fmt.Println(repoPath)
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
			var paths []string
			// Append repositories to check
			if len(args) == 0 {
				for _, p := range configurationFileStructure.VisibleRepositories {
					paths = append(paths, p)
				}
			} else {
				for _, repository := range args {
					repoPath, ok := configurationFileStructure.VisibleRepositories[repository]
					if ok {
						paths = append(paths, repoPath)
					} else {
						traces.WarningTracer.Printf("%s cannot be found in your visible repositories\n", repository)
					}
				}
			}
			var wg sync.WaitGroup
			for _, repository := range paths {
				wg.Add(1)
				go func(repoPath string) {
					defer wg.Done()
					cGitObj := gitManip.New(repoPath)
					cGitObj.Status()
				}(repository)
			}
			wg.Wait()
		},
	}

	/*switchCmd is a subcommand to switch the visibility of the current git repository.
	 */
	// var switchCmd = &cobra.Command{
	// 	Use:   "switch",
	// 	Short: "Switch the visibility of the current git repository (given by the current path)",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		// Get the path where the command has been executed
	// 		currentDir, err := os.Getwd()
	// 		if err != nil {
	// 			log.Fatalln("there was a problem retrieving the current directory")
	// 		}
	// 		if err := configurationFileStructure.Extract(true); err != nil {
	// 			traces.ErrorTracer.Fatalln(err)
	// 		}
	// 		// Recognizing the repository in the configuration file
	// 		repositoryName := filepath.Base(currentDir)
	// 		localRepositories := configurationFileStructure.Repositories
	// 		repositoryPathIndex := utils.SliceIndex(len(localRepositories), func(i int) bool { return localRepositories[i].Name == repositoryName })
	// 		if repositoryPathIndex == -1 {
	// 			log.Fatalf("the repository '%s' does not exists in the configuration file", filepath.Base(currentDir))
	// 			os.Exit(2)
	// 		}
	// 		repositoryGroupIndex := utils.SliceIndex(len(localRepositories[repositoryPathIndex].Paths), func(i int) bool { return localRepositories[repositoryPathIndex].Paths[i].Path == currentDir })
	// 		if repositoryGroupIndex == -1 {
	// 			log.Fatalf("the file path '%s' does not exists in the configuration file", currentDir)
	// 			os.Exit(2)
	// 		}
	// 		// The repository is now recognized
	// 		currentGroupName := utils.GetLocalhost()
	// 		groupIndex := utils.SliceIndex(len(configurationFileStructure.Groups), func(i int) bool { return configurationFileStructure.Groups[i].Name == currentGroupName })
	// 		if groupIndex == -1 {
	// 			log.Fatalln("your localhost is not recognized, please to crawl first")
	// 			os.Exit(2)
	// 		}
	// 		visibleRepositories := configurationFileStructure.VisibleRepositories
	// 		// visibleRepositories := configurationFileStructure.Groups[groupIndex].VisibleRepositories
	// 		repositoryIndex := utils.SliceIndex(len(visibleRepositories), func(i int) bool { return visibleRepositories[i].Name == repositoryName })
	// 		if repositoryIndex == -1 {
	// 			configurationFileStructure.AddRepository(currentDir, consts.VisibleFlag)
	// 			fmt.Printf("The repository %s has been added!\n", repositoryName)
	// 		} else {
	// 			configurationFileStructure.VisibleRepositories = append(visibleRepositories[:repositoryIndex], visibleRepositories[repositoryIndex+1:]...)
	// 			fmt.Printf("The repository %s has been removed!\n", repositoryName)
	// 		}
	// 		fmt.Println(configurationFileStructure.VisibleRepositories)
	// 	},
	// }

	// rootCmd.AddCommand(crawlCmd, pathCmd, stateCmd)

	rootCmd.AddCommand(addCmd, crawlCmd, initCmd, pathCmd, stateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
