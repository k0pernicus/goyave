package walk

import (
	"os"
	"path/filepath"

	"github.com/k0pernicus/goyave/consts"
	"github.com/k0pernicus/goyave/traces"
)

/*RetrieveGitRepositories returns an array of strings, which represent paths to git repositories.
 *Also, this function returns an error type, that is corresponding to the Walk function behaviour (ok or not).
 */
func RetrieveGitRepositories(rootpath string) ([]string, error) {
	var gitPaths []string
	err := filepath.Walk(rootpath, func(pathdir string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() && filepath.Base(pathdir) == consts.GitFileName {
			fileDir := filepath.Dir(pathdir)
			traces.DebugTracer.Printf("Just found in hard drive %s\n", fileDir)
			gitPaths = append(gitPaths, fileDir)
		}
		return nil
	})
	return gitPaths, err
}
