# goyave
A local console-based git projects manager

### General

`goyave` is a small and simple command-line tool to interact with your local git repositories, in order to keep an eye on them.   
It creates and updates a TOML file (in your `$HOME`), to perform the speed-up of interactions.

`goyave` allows you to get some informations about _dirty_ git repositories in your system (a _dirty_ repository is a repository that contains non-commited files, modified files, etc...).   
In order to get updates on repositories you are interested in, `goyave` uses a binary system that consists in telling him what are the repositories you are interested for (in this project, we call them _VISIBLE_ repositories).

### Screenshot

![Simple screenshot](./pictures/goyave.png)

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
    Group = yourLocalHostname

# A list of repositories
# The repository may be visible or hidden, according to the group configuration
[[repositories]]
    Name = repositoryName
    URL  = remoteURL
    [[repositories.local]]
        Path = theLocalPath

[[group]]
    Name = aLocalHostname
    VisibleRepositories = [repo1.Name, repo2.Name, etc...]
```

### How to use it?

* `go get github.com/k0pernicus/goyave`
* `goyave crawl` (needed!)
* `goyave state`

### Troubleshootings

* Please to make sure that the 25th version of [libgit2](https://libgit2.github.com/) is installed on your computer.

### LICENSE

MIT License
