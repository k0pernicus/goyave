package gitManip

import (
	"fmt"

	"github.com/fatih/color"

	git "gopkg.in/libgit2/git2go.v24"
)

/*Map to match the RepositoryState enum type with a string
 */
var repositoryStateToString = map[git.RepositoryState]string{
	git.RepositoryStateNone:                 "None",
	git.RepositoryStateMerge:                "Merge",
	git.RepositoryStateRevert:               "Revert",
	git.RepositoryStateCherrypick:           "Cherrypick",
	git.RepositoryStateBisect:               "Bisect",
	git.RepositoryStateRebase:               "Rebase",
	git.RepositoryStateRebaseInteractive:    "Rebase Interactive",
	git.RepositoryStateRebaseMerge:          "Rebase Merge",
	git.RepositoryStateApplyMailbox:         "Apply Mailbox",
	git.RepositoryStateApplyMailboxOrRebase: "Apply Mailbox or Rebase",
}

var fileStateToString = map[git.Status]string{
	git.StatusIndexNew:     "You forgot to commit a new file!",
	git.StatusIgnored:      "Ignored",
	git.StatusConflicted:   "Conflicted",
	git.StatusWtNew:        "New file in your working tree!",
	git.StatusWtModified:   "Modified file in your working tree!",
	git.StatusWtDeleted:    "Deleted file in your working tree!",
	git.StatusWtTypeChange: "Type change detected in your working tree!",
	git.StatusWtRenamed:    "Renamed file in your working tree!",
}

/*Global variable to set the StatusOption parameter, in order to list each file status
 */
var statusOption = git.StatusOptions{
	Show:     git.StatusShowIndexAndWorkdir,
	Flags:    git.StatusOptIncludeUntracked,
	Pathspec: []string{},
}

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

/*isAccesible returns the information that is the current git repository is existing or not.
 *This method returns a boolean value: true if the git repository is still accesible (still exists), or false if not.
 */
func (g *GitObject) isAccessible() bool {
	return g.accessible == nil
}

/*Status prints the current status of the repository, accessible via the structure path field.
 *This method works only if the repository is accessible.
 */
func (g *GitObject) Status() {
	if g.isAccessible() {
		if err := g.printChanges(); err != nil {
			color.RedString("Impossible to get stats from %s, due to error %s", g.path, err)
		}
	}
}

/*getDiffWithWT returns the difference between the working tree and the index, for the current git repository.
 *If there is an error processing the request, it returns an error.
 */
func (g *GitObject) getDiffWithWT() (*git.Diff, error) {
	// Get the index of the repository
	currentIndex, err := g.repository.Index()
	if err != nil {
		return nil, err
	}
	// Get the default diff options
	defaultDiffOptions, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	// Check the difference between the working directory and the index
	diff, err := g.repository.DiffIndexToWorkdir(currentIndex, &defaultDiffOptions)
	if err != nil {
		return nil, err
	}
	return diff, nil
}

/*printChanges prints out all changes for the current git repository.
 *If there is an error processing the request, it returns this one.
 */
func (g *GitObject) printChanges() error {
	diff, err := g.getDiffWithWT()
	if err != nil {
		return err
	}
	numDeltas, err := diff.NumDeltas()
	if err != nil {
		return err
	}
	if numDeltas > 0 {
		fmt.Printf("%s %s\t[%d modification(s)]\n", color.RedString("/!\\"), g.path, numDeltas)
		for i := 0; i < numDeltas; i++ {
			delta, _ := diff.GetDelta(i)
			currentStatus := delta.Status
			newFile := delta.NewFile.Path
			oldFile := delta.OldFile.Path
			switch currentStatus {
			case git.DeltaAdded:
				fmt.Printf("\t===> %s has been added!\n", newFile)
			case git.DeltaDeleted:
				fmt.Printf("\t===> %s has been deleted!\n", newFile)
			case git.DeltaModified:
				fmt.Printf("\t===> %s has been modified!\n", newFile)
			case git.DeltaRenamed:
				fmt.Printf("\t===> %s has been renamed to %s!\n", oldFile, newFile)
			case git.DeltaUntracked:
				fmt.Printf("\t===> %s is untracked - please to add it or update the gitignore file!\n", newFile)
			case git.DeltaTypeChange:
				fmt.Printf("\t===> the type of %s has been changed from %d to %d!", newFile, delta.OldFile.Mode, delta.NewFile.Mode)
			}
		}
	}
	return nil
}

/*getUntrackedFiles is an old version of a tracking function for Goyave.
 *It is already deprecated.
 */
func (g *GitObject) getUntrackedFiles() {
	fmt.Printf("[%s]...", g.path)
	if untrackedFields, err := g.repository.StatusList(&statusOption); err == nil {
		count, _ := untrackedFields.EntryCount()
		if count == 0 {
			fmt.Println(color.GreenString("\tOK!"))
			return
		}
		fmt.Printf(color.RedString("\t%d untracked files!\n", count))
		for i := 0; i < count; i++ {
			statusEntry, _ := untrackedFields.ByIndex(i)
			fmt.Printf("\t%s\n", fileStateToString[statusEntry.Status])
		}
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
