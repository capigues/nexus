package main

import (
	"fmt"
	"io"

	"github.com/capigues/nexus/cmd/nexus/require"
	"github.com/spf13/cobra"
)

const addDesc = `
ADD MORE
`

type addOptions struct {
	name string
	url  string
}

func newAddCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	o := &addOptions{}

	cmd := &cobra.Command{
		Use:   "add NAME URL",
		Short: "add a model serving api at url",
		Long:  addDesc,
		Args:  require.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return nil, cobra.ShellCompDirectiveDefault
			}

			activeHelpMsg := "This command does not take any more arguments (but may accept flags)."
			return cobra.AppendActiveHelp(nil, activeHelpMsg), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.name = args[0]
			o.url = args[1]

			fmt.Fprintf(out, "Adding %s\n", o.name)
			servers.Add(o.name, o.url)

			return nil
		},
	}

	return cmd
}
