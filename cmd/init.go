// SPDX-FileCopyrightText: 2020 Alvar Penning
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cmd

import (
	"os"
	"os/exec"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var exampleConfig = strings.TrimSpace(`
[github]
# Name of your GitHub organization.
orga = "your-github-orga"

# A personal API access token.
# Needs at least "repo" and "admin:repo_hook -> read:repo_hook" permissions.
#
# https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token
token = "github-api-token"

[branch]
# Branch to pull.
pull = "master"

# Branch to push.
push = "feedback"
`)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init directory",
	Short: "Initialize new directory containing a configuration",
	Long: `Create a new directory for the git repositories.

A dummy configuration will be created and opened with the $EDITOR.`,
	Run:  initAction,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// initAction will be executed for the init command.
func initAction(cmd *cobra.Command, args []string) {
	directory := args[0]
	if stat, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.Mkdir(directory, 0755); err != nil {
			log.WithField("directory", directory).WithError(err).Fatal("Cannot create directory")
		}
	} else if err == nil && stat.IsDir() {
		log.WithField("directory", directory).Debug("Directory does already exists")
	} else {
		log.WithField("directory", directory).WithError(err).Fatal("Unknown error")
	}

	configFile := path.Join(directory, ".github-orga-sync.toml")
	if _, err := os.Stat(configFile); err == nil {
		log.WithField("config", configFile).Info("Configuration does already exists")
		return
	}

	if config, err := os.Create(configFile); err != nil {
		log.WithField("config", configFile).WithError(err).Fatal("Cannot create configuration")
	} else if _, err := config.WriteString(exampleConfig); err != nil {
		log.WithField("config", configFile).WithError(err).Fatal("Cannot write configuration")
	} else if err := config.Close(); err != nil {
		log.WithField("config", configFile).WithError(err).Fatal("Cannot close configuration")
	}

	if editor := os.Getenv("EDITOR"); editor == "" {
		log.WithField("config", configFile).Warn("Cannot find $EDITOR, please edit the configuration manually")
	} else {
		editorCmd := exec.Command(editor, configFile)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			log.WithField("config", configFile).WithError(err).Fatal("Cannot edit configuration")
		}
	}

	log.WithField("directory", directory).Info("Directory is prepared")
}
