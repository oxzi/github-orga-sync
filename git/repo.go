package git

// Repo describes a git repo hosted at GitHub.
type Repo struct {
	Name   string
	SshUrl string

	Ref map[string]string
}
