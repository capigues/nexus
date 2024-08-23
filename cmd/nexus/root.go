package main

import (
	"io"

	"github.com/spf13/cobra"
)

const rootDesc = `
ADD MORE LATER
`

func NewRootCmd(servers *ModelServers, out io.Writer, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "nexus",
		Short: "model serving api management tool",
		Long:  rootDesc,
	}

	cmd.AddCommand(
		newAddCommand(servers, out),
		newRemoveCommand(servers, out),
	// newUpdateCommand(servers, out),
	// newListCommand(servers,out),
	)

	return cmd, nil
}
