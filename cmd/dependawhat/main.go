package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "dependawhat",
		Short: "Check for open Dependabot PRs",
		Long: `A read-only tool to check for open Dependabot pull requests.

Lists all open Dependabot PRs across configured repositories with their
CI status and deny list information. Perfect for monitoring dependency
updates without the ability to approve, recreate, or close PRs.

Configuration can be provided via YAML file or command-line flags.`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dependawhat/config.yaml)")
	rootCmd.PersistentFlags().String("github-token", "", "GitHub token (defaults to USER_GITHUB_TOKEN env var)")
	rootCmd.PersistentFlags().StringSlice("deny-packages", []string{}, "Packages to deny")
	rootCmd.PersistentFlags().StringSlice("deny-orgs", []string{}, "Organizations to deny")

	// Bind flags to viper
	viper.BindPFlag("github-token", rootCmd.PersistentFlags().Lookup("github-token"))
	viper.BindPFlag("deny-packages", rootCmd.PersistentFlags().Lookup("deny-packages"))
	viper.BindPFlag("deny-orgs", rootCmd.PersistentFlags().Lookup("deny-orgs"))

	// Add subcommand
	rootCmd.AddCommand(checkCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search for config in home directory
		viper.AddConfigPath(filepath.Join(home, ".dependawhat"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Bind environment variables
	viper.SetEnvPrefix("DEPENDAWHAT")
	viper.AutomaticEnv()

	// Also check for USER_GITHUB_TOKEN specifically
	viper.BindEnv("github-token", "USER_GITHUB_TOKEN")

	// Read config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
