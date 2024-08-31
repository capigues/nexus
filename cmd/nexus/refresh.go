package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var refreshDesc = `
Get updated info for all APIs managed by Nexus by omitting NAME parameter. Specify a name of saved API to only refresh a specific API.
`

type refreshOptions struct {
	name string
}

func newRefreshCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	o := &refreshOptions{}
	cmd := &cobra.Command{
		Use:   "refresh NAME",
		Short: "optional NAME argument",
		Long:  refreshDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("you can only specify one saved API to refresh at a time")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				for _, server := range *servers {
					err := server.GetInfo()
					if err != nil {
						fmt.Fprintf(out, "%v\n", err.Error())
						return
					}

					servers.Update(server.Name, server)
				}
				fmt.Fprintf(out, "Refreshed all APIs\n")
				return
			}

			o.name = args[0]
			for _, server := range *servers {
				if server.Name == o.name {
					err := server.GetInfo()
					if err != nil {
						fmt.Fprintf(out, "%v\n", err.Error())
						return
					}

					servers.Update(server.Name, server)

					fmt.Fprintf(out, "Refreshed %v API\n", o.name)
					return
				}
			}
			fmt.Fprintf(out, "Could not find API with name: %v\n", o.name)

		},
	}

	return cmd
}
