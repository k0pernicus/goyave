package gitManip

import (
	"fmt"

	"github.com/fatih/color"

	git "gopkg.in/libgit2/git2go.v24"
)

/*GitObject contains informations about the current git repository
 *
 *The structure is:
 *  accessible:
 *		Is the repository still exists in the hard drive?
 *	path:
 *		The path file.
 *	repository:
 *		The object repository.
 */
type GitObject struct {
	accessible error
	path       string
	repository git.Repository
}

/*New is a constructor for GitObject
 *
 * It neeeds:
 *	path:
 *		The path of the current repository.
 */
func New(path string) *GitObject {
	r, err := git.OpenRepository(path)
	return &GitObject{accessible: err, path: path, repository: *r}
}

func (g *GitObject) isAccessible() bool {
	return g.accessible == nil
}

/*Status prints the current status of the repository, accessible via the structure path field.
 *This method works only if the repository is accessible.
 */
func (g *GitObject) Status() {
	if g.isAccessible() {
		fmt.Printf("The status of %s is: %s\n", g.path, g.repository.State())
	}
}

/*List lists the path and the accessibility of a list of git repositories
 */
func List(repositories *[]GitObject) {
	for _, object := range *repositories {
		fmt.Printf("* %s ", object.path)
		if object.isAccessible() {
			fmt.Println(color.GreenString(" [accessible]"))
		} else {
			fmt.Println(color.RedString(" [not accessible]"))
		}
	}
}
