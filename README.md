<!--
SPDX-FileCopyrightText: 2020 Alvar Penning

SPDX-License-Identifier: GPL-3.0-or-later
-->

# github-orga-sync

`github-orga-sync` is a simple tool to synchronize all repositories from a GitHub organization.

The intended workflow is to deal with a GitHub Classroom "only" organization.
New student repositories will be cloned or updated based on their master branch.
A feedback can be pushed from a feedback branch afterwards.

> :warning: __This software is quite new and untested! Don't rely on it. Here be dragons.__


## Installation

Install a current [Go](https://golang.org/) version.

```
git clone https://github.com/oxzi/github-orga-sync
cd github-orga-sync
go build
# Put the github-orga-sync binary into your $PATH
```


## Workflow

`github-orga-sync`'s workflow is very similar to that of `git` itself.
At first, one initializes the main directory to contain all repositories.
Afterwards, one pulls all repositories.
After some changes, one pushes back those.
The last two steps will be repeated.

### 1. Initialize directory

First, we create a new directory `repos`.
The configuration file will be opened with our `$EDITOR`.

```
$ github-orga-sync init repos
> INFO Directory is prepared                         directory=repos

$ cd repos
```

### 2. Pull

The internal `git clone` will be used with the SSH URI.
To avoid having to enter the passphrase every time, it is recommended to set up `ssh-agent`.

```
$ github-orga-sync pull
> INFO Fetching repositories from GitHub..           branch=master organization=my-github-orga
> INFO There are 3 repositories in total
> INFO Created new repository                        repository=ubung-01-testteam
> INFO Created new repository                        repository=ubung-01-team23
> INFO Created new repository                        repository=ubung-01-team42
> INFO Finished pull successfully                    created=3 updated=0
```

### 3. Push

We make some changes in our `branch.push` branch (defaults to `feedback`).
Afterwards, we push them back.

For this example, we only alter one repository, `ubung-01-testteam`.

```
$ github-orga-sync push
> INFO Fetching repositories from GitHub..           branch=master organization=my-github-orga
> INFO There are 3 repositories in total
> INFO Pushed updates to remote                      repository=ubung-01-testteam
> WARN Local repository does not contain branch      branch=feedback repository=ubung-01-team23
> WARN Local repository does not contain branch      branch=feedback repository=ubung-01-team42
> INFO Finished push successfully                    updated=1
```


## License

GNU GPLv3 or later.
