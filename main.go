package main

import (
	"fmt"
	"os"
	"swgw/command/commands"

	"github.com/spf13/cobra"
)

func main() {

	cmd := newSwgwCommand()

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}

}

func newSwgwCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                    "swgw [OPTIONS] COMMAND [ARGS...]",
		Short:                  "create cloudformation yml",
		SilenceErrors:          true,
		SilenceUsage:           true,
		BashCompletionFunction: "cf",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return fmt.Errorf("swgw: '%s' is not a swgw command.\nSee 'swgw --help'", args[0])
		},
		Version:               fmt.Sprintf("version:%s\n", "1.0"),
		DisableFlagsInUseLine: true,
	}

	cmd.SetOutput(os.Stdout)
	cmd.PersistentFlags().StringP("file", "f", "cloudformation.yml", "-f ecs.yml")
	commands.AddCommands(cmd)
	return cmd

}
