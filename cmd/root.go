/*
Package cmd contains all the commands that are available in the test-inspector CLI.
Copyright © 2022 Egor Romanov egor@supabase.io

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	host        string
	user        string
	password    string
	resultsPath string
	reportType  string
	versionID   int32
)

var (
	// Version is used to show the version of the CLI build.
	Version = "dev"
	// CommitHash is used to show the commit hash of the CLI build.
	CommitHash = "n/a"
	// BuildTime is used to show the time of the CLI build.
	BuildTime = "n/a"
	// SupabaseKey is used to as a apiKey for the Supabase API.
	SupabaseKey = ""
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "test-inspector",
	Short: "An app to inspect your test suite and upload results to test-inspector",
	Long: `
An app to inspect your test suite and upload results to test-inspector.
Version: ` + Version + `
Commit: ` + CommitHash + `
Build time: ` + BuildTime + `
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "",
		"config file (default is $HOME/.test-inspector.yaml)")
	rootCmd.PersistentFlags().StringVarP(
		&host, "host", "H", "https://gryakvuryfsrgjohzhbq.supabase.co",
		"url for test-inspector backend")
	rootCmd.PersistentFlags().StringVarP(
		&user, "user", "u", "",
		"test-inspector user email")
	rootCmd.PersistentFlags().StringVarP(
		&password, "password", "w", "",
		"test-inspector user password")
	rootCmd.PersistentFlags().StringVarP(
		&resultsPath, "resultsPath", "f", "./allure-results",
		"path to the directory with allure results")
	rootCmd.PersistentFlags().Int32VarP(
		&versionID, "versionID", "v", 0, "version ID in test-inspector (required)")
	rootCmd.PersistentFlags().StringVarP(
		&reportType, "type", "t", "allure", "report type (possible values: allure, junit)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("resultsPath", rootCmd.PersistentFlags().Lookup("resultsPath"))
	viper.BindPFlag("versionID", rootCmd.PersistentFlags().Lookup("versionID"))
	viper.BindPFlag("type", rootCmd.PersistentFlags().Lookup("type"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".test-inspector" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".test-inspector")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func validateFlags() error {
	if user == "" {
		return fmt.Errorf("user email is required")
	}
	if password == "" {
		return fmt.Errorf("user password is required")
	}
	if versionID == 0 {
		return fmt.Errorf("versionID is required")
	}
	return nil
}

func validateVersionID() error {
	if versionID == 0 {
		return fmt.Errorf("versionID is required")
	}
	return nil
}
