package commands

import (
	"swgw/command"
	"swgw/command/cf"
	"swgw/command/cost"

	"github.com/spf13/cobra"
)

//AddCommands コマンドを追加する
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(cf.NewCFCommand())
	cmd.AddCommand(command.NewPLCommand())
	cmd.AddCommand(cost.NewCostCommand())
}
