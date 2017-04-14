/*Package enums implements enums for the entire project
 */
package enums

// VISIBILITY is a string alias to get the local visibility of git repositories
type VISIBILITY string

/*Enumeration to specify the local visibility of git repositories
 *
 *Those enumerations are:
 *	VISIBLE:
 *		Visible repository - can be access through the program
 *	HIDDEN:
 *		An ignored repository - canno't be access through the program
 */
const (
	VISIBLE = VISIBILITY("VISIBLE")
	HIDDEN  = VISIBILITY("HIDDEN")
)
