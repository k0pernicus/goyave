# goyave
A local console-based git projects manager

### General

`goyave` is a small and simple command-line tool to interact with your local git repositories, in order to keep an eye on them.   
It creates and updates a TOML file (in your `$HOME`), to perform the speed-up of interactions.

`goyave` allows you to get some informations about _dirty_ git repositories in your system (a _dirty_ repository is a repository that contains non-commited files, modified files, etc...).   
In order to get updates on repositories you are interested in, `goyave` uses a binary system that consists in telling him what are the repositories you are interested for (in this project, we call them _VISIBLE_ repositories).

### The configuration file

The configuration file is available at `$HOME/.goyave`.

This is the structure of this configuration file:

```TOML
# The name of the user - typically the name of your session account
Author = yourSessionAccount

# Some local informations
[local]
    # The target to store new git repositories
    # Available properties are:
    # * "VISIBLE" - Store each new git repository as a _VISIBLE_ repository
    # * "HIDDEN" - Store each new git repository as a _HIDDEN_ one
    DefaultTarget = "VISIBLE"
    # The group you are using to perform some actions - typically, your hostname
    Group = yourHostname

# A list of visible repositories
# A visible repository is a repository you want some updates on
[[visible]]
    Name = repositoryName
    Path = repositoryPath

# A list of hidden repositories
# An hidden repository is a repository you do not want updates on
[[hidden]]
    Name = repositoryName
    Path = repositoryPath
```

### How to use it?

* `go get github.com/k0pernicus/goyave`
* `goyave crawl` (needed!)
* `goyave state`

### LICENSE

MIT License