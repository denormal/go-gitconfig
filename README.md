# go-gitconfig

Package `gitconfig` provides an interface to git configuration properties
as returned by `"git config --list`. `gitconfig` provides access to local,
global and system configuration, as well as the effective configuration
for the given git working copy. `gitconfig` attempts to use the locally
installed `git` executable via
[go-gittools](https://github.com/denormal/go-gittools).

See [git-config](https://git-scm.com/docs/git-config) for more information.

```go
import "github.com/denormal/go-gitconfig"

// load the git configuration for a particular repository
config, err := gitconfig.NewWithPath("/my/git/working/copy")
if err != nil {
    panic(err)
}

// extract the core.* properties
core := config.Find("core.*")

// extract the git user's name
user := config.Get("user.name")
```

For more information see `godoc github.com/denormal/go-gitconfig`.

## Installation

`go-gitconfig` can be installed using the standard Go approach:

```go
go get github.com/denormal/go-gitconfig
```

## License

Copyright (c) 2016 Denormal Limited

[MIT License](LICENSE)
