package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/oxzi/github-orga-sync/git"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push local changes to existing upstream repositories",
	Run:   pushAction,
}

func init() {
	rootCmd.AddCommand(pushCmd)
}

// pushAction will be executed for the push command.
func pushAction(cmd *cobra.Command, args []string) {
	parseConfig()

	log.WithFields(log.Fields{
		"organization": viper.GetString("github.orga"),
		"branch":       viper.GetString("branch.push"),
	}).Info("Fetching repositories from GitHub..")

	repos, err := git.FetchGitHubRepos(
		viper.GetString("github.orga"),
		"refs/heads/"+viper.GetString("branch.push"),
		viper.GetString("github.token"))
	if err != nil {
		log.WithError(err).Fatal("Fetching failed")
	}

	log.Infof("There are %d repositories in total", len(repos))
	var countUpdated int

	for _, repo := range repos {
		updated, err := repo.Push(viper.GetString("branch.push"))

		if err != nil {
			log.WithField("repository", repo).WithError(err).Fatal("Updating failed")
		} else if updated {
			log.WithField("repository", repo).Info("Pushed updates to remote")
			countUpdated++
		} else {
			log.WithField("repository", repo).Debug("No update for this repository")
		}
	}

	log.WithField("updated", countUpdated).Info("Finished push successfully")
}
