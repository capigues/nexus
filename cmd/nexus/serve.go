package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const serveDesc = `
Serve all APIs saved in Nexus at one single endpoint
`

func newServeCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "start nexus proxy server",
		Long:  serveDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := servers.Serve()
			if err != nil {
				fmt.Fprintf(out, "%v\n", err.Error())
			}

			fmt.Fprintf(out, "")
		},
	}

	return cmd
}
