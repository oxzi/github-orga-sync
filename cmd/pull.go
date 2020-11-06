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
	}).Info("Feching repositories from GitHub..")

	repos, err := git.FetchGitHubRepos(
		viper.GetString("github.orga"),
		"refs/heads/"+viper.GetString("branch.pull"),
		viper.GetString("github.token"))
	if err != nil {
		log.WithError(err).Fatal("Fetching errored")
	}

	log.Infof("There are %d repositories in total", len(repos))
}
