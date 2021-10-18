package cmd

import (
	"github.com/murex/tcr-cli/cli"
	"github.com/murex/tcr-engine/runmode"
	"github.com/murex/tcr-engine/settings"
	"github.com/spf13/cobra"
)

// mobCmd represents the mob command
var mobCmd = &cobra.Command{
	Use:   "mob",
	Short: "Run TCR in mob mode",
	Long: `
When used in "mob" mode, TCR ensures that any commit
is shared with other participants through calling git push-pull.

This subcommand runs directly in the terminal (no GUI).`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Mode = runmode.Mob{}
		params.AutoPush = params.Mode.AutoPushDefault()
		params.PollingPeriod = settings.DefaultPollingPeriod
		u := cli.New(params)
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(mobCmd)
}
