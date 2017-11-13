package gitManip

import (
	"fmt"

	"github.com/fatih/color"

	"bytes"

	"github.com/k0pernicus/goyave/traces"
	git "gopkg.in/libgit2/git2go.v26"
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

/*Clone is cloning a given repository, from a public URL
 *
 * It needs:
 * path:
 *		The local path to clone the repository.
 *	URL:
 *		The remote URL to fetch the repository.
 */
func Clone(path, URL string) error {
	_, err := git.Clone(URL, path, &git.CloneOptions{})
	return err
}

/*GetRemoteURL returns the associated remote URL of a given local path repository
 *
 * It needs:
 *	path
 *		The local path of a git repository
 */
func GetRemoteURL(path string) string {
	r, err := git.OpenRepository(path)
	if err != nil {
		fmt.Println("The repository can't be opened")
		return ""
	}
	remoteCollection := r.Remotes
	originRemote, err := remoteCollection.Lookup("origin")
	if err != nil {
		traces.WarningTracer.Printf("can't lookup origin remote URL for %s", path)
		return ""
	}
	return originRemote.Url()
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
	} else {
		color.RedString("Repository %s not found!", g.path)
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
	// Get the default diff options, and add it custom flags
	defaultDiffOptions, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	defaultDiffOptions.Flags = defaultDiffOptions.Flags | git.DiffIncludeUntracked | git.DiffIncludeTypeChange
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
	var buffer bytes.Buffer
	if err != nil {
		return err
	}
	numDeltas, err := diff.NumDeltas()
	if err != nil {
		return err
	}
	headDetached, err := g.repository.IsHeadDetached()
	if err != nil {
		return err
	}
	if headDetached {
		outputHead := fmt.Sprintf("%s", color.RedString("\t/!\\ The repository's HEAD is detached! /!\\\n"))
		buffer.WriteString(outputHead)
	}
	if numDeltas > 0 {
		buffer.WriteString(fmt.Sprintf("%s %s\t[%d modification(s)]\n", color.RedString("✘"), g.path, numDeltas))
		for i := 0; i < numDeltas; i++ {
			delta, _ := diff.GetDelta(i)
			currentStatus := delta.Status
			newFile := delta.NewFile.Path
			oldFile := delta.OldFile.Path
			switch currentStatus {
			case git.DeltaAdded:
				buffer.WriteString(fmt.Sprintf("\t===> %s has been added!\n", color.MagentaString(newFile)))
			case git.DeltaDeleted:
				buffer.WriteString(fmt.Sprintf("\t===> %s has been deleted!\n", color.MagentaString(newFile)))
			case git.DeltaModified:
				buffer.WriteString(fmt.Sprintf("\t===> %s has been modified!\n", color.MagentaString(newFile)))
			case git.DeltaRenamed:
				buffer.WriteString(fmt.Sprintf("\t===> %s has been renamed to %s!\n", color.MagentaString(oldFile), color.MagentaString(newFile)))
			case git.DeltaUntracked:
				buffer.WriteString(fmt.Sprintf("\t===> %s is untracked - please to add it or update the gitignore file!\n", color.MagentaString(newFile)))
			case git.DeltaTypeChange:
				buffer.WriteString(fmt.Sprintf("\t===> the type of %s has been changed from %d to %d!", color.MagentaString(newFile), delta.OldFile.Mode, delta.NewFile.Mode))
			}
		}
	} else {
		buffer.WriteString(fmt.Sprintf("%s %s\n", color.GreenString("✔"), g.path))
	}
	repository_head, err := g.repository.Head()
	if err == nil {
		repository_id := repository_head.Target()
		commits_ahead, _, err := g.repository.AheadBehind(repository_id, repository_id)
		if err != nil {
			buffer.WriteString(fmt.Sprintf("%s", color.RedString("\tAn error occured checking the ahead/behind commits...\n")))
		} else if commits_ahead != 0 {
			buffer.WriteString(fmt.Sprintf("\tYou need to push the last modifications!\n"))
		}
	}
	fmt.Print(buffer.String())
	return nil
}
