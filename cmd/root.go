package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var logDebug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-orga-sync",
	Short: "Bulk pull or push repositories from a GitHub organization",
	Long: `Simple tool to synchronize all repositories from a GitHub organization.

The intended workflow is to deal with a GitHub Classroom "only" organization.
New student repositories will be cloned or updated based on their master branch.
A feedback can be pushed from a feedback branch afterwards.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogging)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".github-orga-sync.toml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&logDebug, "verbose", "v", false, "more verbosity, debug logging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()
}

// initLogging configures the logging level.
func initLogging() {
	if logDebug {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
}

// parseConfig or exit trying.
func parseConfig() {
	errorf := func(format string, args ...interface{}) {
		_, _ = fmt.Fprintf(os.Stderr, format, args...)
		os.Exit(1)
	}

	if err := viper.ReadInConfig(); err != nil {
		errorf("Cannot parse configuration: %v\n", err)
	}

	for _, key := range []string{"github.orga", "github.token", "branch.pull", "branch.push"} {
		if !viper.IsSet(key) {
			errorf("Configuration key \"%s\" is missing\n", key)
		}
	}
}
