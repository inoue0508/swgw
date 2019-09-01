package commands

import (
	"swgw/command"
	"swgw/command/cf"

	"github.com/spf13/cobra"
)

//AddCommands コマンドを追加する
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(cf.NewCFCommand())
	cmd.AddCommand(command.NewPLCommand())
}
