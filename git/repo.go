package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Repo describes a git repo hosted at GitHub.
type Repo struct {
	Name   string
	SshUrl string

	Rev map[string]string
}

func (repo Repo) String() string {
	return repo.Name
}

// Pull this Repo from GitHub. If it does not exists, clone it.
func (repo Repo) Pull(branch string) (created, updated bool, err error) {
	if stat, statErr := os.Stat(repo.Name); os.IsNotExist(statErr) {
		created = true
		updated = true
		err = repo.clone()
	} else if stat.IsDir() {
		updated, err = repo.pull(branch)
	} else {
		err = fmt.Errorf("%s is not a directory", repo.Name)
	}
	return
}

// clone this Repo.
func (repo Repo) clone() error {
	gitCmd, err := exec.LookPath("git")
	if err != nil {
		return err
	}

	cmd := exec.Command(gitCmd, "clone", "-q", repo.SshUrl)
	log.WithFields(log.Fields{
		"repository": repo,
		"cmd":        cmd,
	}).Debug("Cloning repository")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.WithField("repository", repo).WithError(err).Errorf("Cloning failed\n%s", output)
	}

	return err
}

// pull this Repo.
func (repo Repo) pull(branch string) (updated bool, err error) {
	gitCmd, gitCmdErr := exec.LookPath("git")
	if gitCmdErr != nil {
		err = gitCmdErr
		return
	}

	checkoutCmd := exec.Command(gitCmd, "-C", repo.Name, "checkout", "-q", branch)
	if output, checkoutErr := checkoutCmd.CombinedOutput(); checkoutErr != nil {
		log.WithField("repository", repo).WithError(checkoutErr).Errorf("Local checkout failed\n%s", output)
		err = checkoutErr
		return
	}

	revParseCmd := exec.Command(gitCmd, "-C", repo.Name, "rev-parse", "HEAD")
	localRevOutput, revParseErr := revParseCmd.CombinedOutput()
	if revParseErr != nil {
		log.WithField("repository", repo).WithError(revParseErr).Error("git rev-parse failed")
		err = revParseErr
		return
	}
	localRev := strings.TrimSpace(string(localRevOutput))
	log.WithFields(log.Fields{
		"repository": repo,
		"branch":     branch,
		"local-rev":  localRev,
	}).Debug("Fetched local rev")

	remoteRev, isRemoteRev := repo.Rev[branch]
	if !isRemoteRev {
		log.WithField("repository", "repo").Warn("Misses remote branch revision")
	} else if localRev == remoteRev {
		log.WithField("repository", "repo").Debug("No new commits, skipping")
		return
	}

	log.WithField("repository", "repo").Debug("Fetching repository")

	pullCmd := exec.Command(gitCmd, "-C", repo.Name, "pull", "-q", "origin", branch)
	if pullOutput, pullErr := pullCmd.CombinedOutput(); pullErr != nil {
		log.WithField("repository", repo).WithError(pullErr).Errorf("Pulling failed\n%s", pullOutput)
		err = pullErr
		return
	}

	updated = true
	return
}

// Push a branch back to GitHub.
func (repo Repo) Push(branch string) (updated bool, err error) {
	if stat, statErr := os.Stat(repo.Name); os.IsNotExist(statErr) {
		log.WithField("repository", repo).Warn("Remote repository does not exist")
		return
	} else if statErr != nil {
		err = statErr
		return
	} else if !stat.IsDir() {
		err = fmt.Errorf("%s is not a directory", repo.Name)
		return
	}

	gitCmd, gitCmdErr := exec.LookPath("git")
	if gitCmdErr != nil {
		err = gitCmdErr
		return
	}

	revParseCmd := exec.Command(gitCmd, "-C", repo.Name, "rev-parse", branch)
	localRevOutput, revParseErr := revParseCmd.CombinedOutput()
	if revParseErr != nil {
		log.WithField("repository", repo).WithError(revParseErr).Error("git rev-parse failed")
		err = revParseErr
		return
	}
	localRev := strings.TrimSpace(string(localRevOutput))
	log.WithFields(log.Fields{
		"repository": repo,
		"branch":     branch,
		"local-rev":  localRev,
	}).Debug("Fetched local rev")

	remoteRev, isRemoteRev := repo.Rev[branch]
	if !isRemoteRev {
		log.WithField("repository", "repo").Warn("Misses remote branch revision")
	} else if localRev == remoteRev {
		log.WithField("repository", "repo").Debug("No new commits, skipping")
		return
	}

	pushCmd := exec.Command(gitCmd, "-C", repo.Name, "push", "-q", "origin", branch)
	if pushOutput, pushErr := pushCmd.CombinedOutput(); pushErr != nil {
		log.WithField("repository", repo).WithError(pushErr).Errorf("Pushing failed\n%s", pushOutput)
		err = pushErr
		return
	}

	updated = true
	return
}
