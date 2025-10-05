package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/johnmayou/git-contributions-cli/internal/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	scanRunClearFlag = "clear"
	scanIgnore       = []string{"venv", ".venv", "node_modules"}
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "manage scanning and caching of git repositories",
}

var scanRunCmd = &cobra.Command{
	Use:   "run [directory]",
	Short: "scan a directory for git repositories",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return scan.Scan(args[0], scanIgnore, scanCacheFile(), viper.GetBool(scanRunClearFlag))
	},
}

var scanShowCmd = &cobra.Command{
	Use:   "show",
	Short: "print cached git repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		repos, err := scan.Load(scanCacheFile())
		if err != nil {
			return err
		}
		for _, repo := range repos {
			fmt.Println(repo)
		}
		return nil
	},
}

func init() {
	scanRunCmd.Flags().Bool(scanRunClearFlag, false, "override saved git repositories")
	viper.BindPFlag(scanRunClearFlag, scanRunCmd.Flags().Lookup(scanRunClearFlag))
	scanCmd.AddCommand(scanRunCmd)

	scanCmd.AddCommand(scanShowCmd)

	rootCmd.AddCommand(scanCmd)
}

func scanCacheFile() string {
	return filepath.Join(viper.GetString(appDirFlag), "scanned.txt")
}
