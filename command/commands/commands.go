package commands

import (
	"swgw/command"
	"swgw/command/cost"

	"github.com/spf13/cobra"
)

//AddCommands コマンドを追加する
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(command.NewPLCommand())
	cmd.AddCommand(cost.NewCostCommand())
}
