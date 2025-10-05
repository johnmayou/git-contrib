package cmd

import (
	"os"
	"time"

	"github.com/johnmayou/git-contrib/internal/scan"
	"github.com/johnmayou/git-contrib/internal/stats"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	statsEmailFlag = "email"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "print a commit heatmap from previously scanned repositories",
	Long: `Reads cached git repository paths saved by the scan command and
prints a GitHub-style commit heatmap based on daily commit volume.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repos, err := scan.Load(scanCacheFile())
		if err != nil {
			return err
		}

		return stats.Stats(
			repos,
			viper.GetString(statsEmailFlag),
			time.Now,
			stats.SnappedToSundayMorning(time.Now().AddDate(-1, 0, 0)),
			os.Stdout,
		)
	},
}

func init() {
	statsCmd.Flags().String(statsEmailFlag, "", "email to filter commits by")
	statsCmd.MarkFlagRequired(statsEmailFlag)
	viper.BindPFlag(statsEmailFlag, statsCmd.Flags().Lookup(statsEmailFlag))

	rootCmd.AddCommand(statsCmd)
}
