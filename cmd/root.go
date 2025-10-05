package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appDirFlag = "app-dir"
)

var rootCmd = &cobra.Command{
	Use:   "git-contrib",
	Short: "visualize your git commit history across multiple repositories",
	Long: `CLI tool that scans for local git repositories and prints a GitHub-style
contribution heatmap based on daily commit volume.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := os.MkdirAll(viper.GetString(appDirFlag), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating app directory: %w", err)
		}
		return nil
	},
}

var userHomeDir = func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("fatal: count not determine user home directory: " + err.Error())
	}
	return home
}()

var (
	defaultAppDir = filepath.Join(userHomeDir, ".git-contrib")
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String(appDirFlag, defaultAppDir, "working app directory")
	viper.BindPFlag(appDirFlag, rootCmd.PersistentFlags().Lookup(appDirFlag))
}

func initConfig() {
	viper.SetEnvPrefix("GIT_CONTRIB")
	viper.AutomaticEnv()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
