package cmd

import (
	"github.com/mengdaming/tcr-cli/cli"
	"github.com/mengdaming/tcr-engine/engine"
	"github.com/mengdaming/tcr-engine/runmode"

	"github.com/spf13/cobra"
)

// soloCmd represents the solo command
var soloCmd = &cobra.Command{
	Use:   "solo",
	Short: "Run TCR in solo mode",
	Long: `
When used in "solo" mode, TCR only commits changes locally.
It never pushes or pulls to a remote repository.

This subcommand runs directly in the terminal (no GUI).`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Mode = runmode.Solo{}
		params.AutoPush = params.Mode.AutoPushDefault()
		params.PollingPeriod = engine.DefaultPollingPeriod
		u := cli.New(params)
		u.Start()
	},
}

func init() {
	rootCmd.AddCommand(soloCmd)
}