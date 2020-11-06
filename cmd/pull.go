package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/oxzi/github-orga-sync/git"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull or clone all new repositories",
	Run:   pullAction,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

// pullAction will be executed for the pull command.
func pullAction(cmd *cobra.Command, args []string) {
	parseConfig()

	log.WithFields(log.Fields{
		"organization": viper.GetString("github.orga"),
		"branch":       viper.GetString("branch.pull"),
	}).Info("Fetching repositories from GitHub..")

	repos, err := git.FetchGitHubRepos(
		viper.GetString("github.orga"),
		"refs/heads/"+viper.GetString("branch.pull"),
		viper.GetString("github.token"))
	if err != nil {
		log.WithError(err).Fatal("Fetching failed")
	}

	log.Infof("There are %d repositories in total", len(repos))
	var countCreated, countUpdated int

	for _, repo := range repos {
		created, updated, err := repo.Pull(viper.GetString("branch.pull"))

		if err != nil {
			log.WithField("repository", repo).WithError(err).Fatal("Updating failed")
		} else if created {
			log.WithField("repository", repo).Info("Created new repository")
			countCreated++
		} else if updated {
			log.WithField("repository", repo).Info("Fetched update for repository")
			countUpdated++
		} else {
			log.WithField("repository", repo).Debug("No update for this repository")
		}
	}

	log.WithFields(log.Fields{
		"created": countCreated,
		"updated": countUpdated,
	}).Info("Finished pull successfully")
}
