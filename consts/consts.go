/*Package consts implements constants for the entire project
 */
package consts

const DefaultUserName = "Thor"

// VisibleFlag is the constant given for a visible repository
const VisibleFlag = "VISIBLE"

// HiddenFlag is the constant given for an hidden repository
const HiddenFlag = "HIDDEN"

// ConfigurationFileName is the configuration file name of Goyave
const ConfigurationFileName = ".goyave"

// GitFileName is the name of the git directory, in a git repository
const GitFileName = ".git"

////////////////
//// ERRORS ////
////////////////

// RepositoryAlreadyExists is an error that raises when an existing path repository is in a list
const RepositoryAlreadyExists = "REPOSITORY_ALREADY_EXISTS"

// ItemIsNotIsSlice is an error that raises when an searched item is not in the given slice
const ItemIsNotInSlice = "ITEM_IS_NOT_IN_SLICE"
