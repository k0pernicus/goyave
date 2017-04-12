package configurationFile

import "testing"
import "bytes"
import "fmt"

func TestDecode(t *testing.T) {
	const configurationExample = `
		author = 'Antonin'

		[local]
		group = 'custom'
		
		[[visible]]
		name = 'visible'
		path = '/home/user/visible'

		[[hidden]]
		name = 'hidden'
		path = '/home/user/hidden'

		[[hidden]]
		name = 'hidden2'
		path = '/home/user/hidden2'

		[[group]]
		name = 'custom'
		repositories = ['path']
	`
	var configurationStructure ConfigurationFile
	Decode(configurationExample, &configurationStructure)
	// Check the Author entry
	if configurationStructure.Author != "Antonin" {
		t.Errorf("The author in the configuration file example is not 'Antonin' but %s.", configurationStructure.Author)
	}
	// Check the Local structure
	if configurationStructure.Local.Group != "custom" {
		t.Errorf("The local group in the configuration file example is not 'custom' but %s.", configurationStructure.Local.Group)
	}
	// Check the first VisibleRepositories structure
	if len(configurationStructure.VisibleRepositories) != 1 {
		t.Errorf("The number of visible git repositories is not good, got %d instead of %d.", len(configurationStructure.VisibleRepositories), 1)
	}
	if configurationStructure.VisibleRepositories[0].Name != "visible" {
		t.Errorf("The name of the first visible entry is not correct, got %s instead of %s.", configurationStructure.VisibleRepositories[0].Name, "visible")
	}
	if configurationStructure.VisibleRepositories[0].Path != "/home/user/visible" {
		t.Errorf("The path of the first visible entry is not correct, got %s instead of %s.", configurationStructure.VisibleRepositories[0].Path, "/home/user/visible")
	}
	// Check the first HiddenRepositories structure
	if len(configurationStructure.HiddenRepositories) != 2 {
		t.Errorf("The number of hidden git repositories is not good, got %d instead of %d.", len(configurationStructure.HiddenRepositories), 1)
	}
	if configurationStructure.HiddenRepositories[0].Name != "hidden" {
		t.Errorf("The name of the first hidden entry is not correct, got %s instead of %s.", configurationStructure.HiddenRepositories[0].Name, "hidden")
	}
	if configurationStructure.HiddenRepositories[0].Path != "/home/user/hidden" {
		t.Errorf("The path of the first hidden entry is not correct, got %s instead of %s.", configurationStructure.HiddenRepositories[0].Path, "/home/user/hidden")
	}
	// Check the first Groups structure
	if len(configurationStructure.Groups) != 1 {
		t.Errorf("The number of groups is not good, got %d instead of %d.", len(configurationStructure.Groups), 1)
	}
	if configurationStructure.Groups[0].Name != "custom" {
		t.Errorf("The name of the first group entry is not correct, got %s instead of %s.", configurationStructure.Groups[0].Name, "custom")
	}
	if len(configurationStructure.Groups[0].Repositories) != 1 {
		t.Errorf("The number of repositories in the first group is not good, got %d instead of %d.", len(configurationStructure.Groups[0].Repositories), 1)
	}
}

func TestEncode(t *testing.T) {
	localStructure := ConfigurationFile{
		Author: "Antonin",
		Local: LocalInformations{
			Group: "local",
		},
		VisibleRepositories: []GitRepository{
			GitRepository{
				Name: "visible_example",
				Path: "/home/user/mypath",
			},
		},
	}
	buffer := new(bytes.Buffer)
	Encode(&localStructure, buffer)
	fmt.Println(buffer)
}
