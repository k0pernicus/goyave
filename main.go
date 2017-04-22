package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"

	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/k0pernicus/goyave/configurationFile"
	"github.com/k0pernicus/goyave/consts"
	"github.com/k0pernicus/goyave/traces"
	"github.com/k0pernicus/goyave/utils"
	"github.com/k0pernicus/goyave/walk"
)

/*getConfigurationFileContent get the content of the local configuration file.
 *If no configuration file has been found, create a default one and set the bytes array.
 */
func getConfigurationFileContent(filePointer *os.File, bytesArray *[]byte) {
	fileState, err := filePointer.Stat()
	// If the file is empty, get the default structure and save it
	if err != nil || fileState.Size() == 0 {
		traces.WarningTracer.Println("No (or empty) configuration file - creating default one...")
		var fileBuffer bytes.Buffer
		defaultStructure := configurationFile.Default()
		defaultStructure.Encode(&fileBuffer)
		*bytesArray = fileBuffer.Bytes()
	} else {
		b, _ := ioutil.ReadAll(filePointer)
		*bytesArray = b
	}
}

func main() {

	// Initialize all different traces structures
	traces.InitTraces(os.Stdout, os.Stderr, os.Stdout, os.Stdout)

	// Get the user home directory
	userHomeDir := utils.GetUserHomeDir()
	if len(userHomeDir) == 0 {
		log.Fatalf("Cannot get the user home dir.\n")
	}
	// Set the configuration path file
	configurationFilePath := path.Join(userHomeDir, consts.ConfigurationFileName)
	filePointer, err := os.OpenFile(configurationFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Cannot open the file %s, due to error '%s'.\n", configurationFilePath, err)
	}
	defer filePointer.Close()
	var bytesArray []byte
	// Get the content of the goyave configuration file
	getConfigurationFileContent(filePointer, &bytesArray)
	var configurationFileStructure configurationFile.ConfigurationFile
	if _, err = toml.Decode(string(bytesArray[:]), &configurationFileStructure); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(configurationFileStructure)
	// Get all git paths, and display them
	gitPaths, err := walk.RetrieveGitRepositories(userHomeDir)
	if err != nil {
		log.Fatalf("There was an error retrieving your git repositories: '%s'\n", err)
	}
	for _, gitPath := range gitPaths {
		if err := configurationFileStructure.AddRepository(gitPath, configurationFileStructure.Local.DefaultTarget); err != nil {
			traces.WarningTracer.Printf("[%s] %s", gitPath, err)
		}
	}
	var outputBuffer bytes.Buffer
	if err := configurationFileStructure.Encode(&outputBuffer); err != nil {
		log.Fatalln("Cannot save the current configurationFile structure!")
	}
	if err := ioutil.WriteFile(configurationFilePath, outputBuffer.Bytes(), 0777); err != nil {
		log.Fatalln("Cannot access to your file to save the configurationFile structure!")
	}
}
