package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const removeDesc = `
Remove API from being managed by Nexus
`

type removeOptions struct {
	name string
}

func newRemoveCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	o := &removeOptions{}

	cmd := &cobra.Command{
		Use:   "remove NAME",
		Short: "remove a model serving api at url",
		Long:  removeDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1); err != nil {
				if len(args) == 0 {
					return errors.New("you must specify the name of the model serving api you are removing")
				}
				if len(args) != 1 {
					return errors.New("this command does not take any more arguments")
				}
			}

			return nil
		},
		// ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// 	if len(args) == 0 {
		// 		return nil, cobra.ShellCompDirectiveDefault
		// 	}

		// 	activeHelpMsg := "This command does not take any more arguments (but may accept flags)."
		// 	return cobra.AppendActiveHelp(nil, activeHelpMsg), cobra.ShellCompDirectiveNoFileComp
		// },
		Run: func(cmd *cobra.Command, args []string) {
			o.name = args[0]

			if err := servers.Remove(o.name); err != nil {
				fmt.Fprintf(out, "Could not remove %s\nnexus: %v", o.name, err.Error())
				return
			}

			fmt.Fprintf(out, "Removed %s\n", o.name)
		},
	}

	return cmd
}
