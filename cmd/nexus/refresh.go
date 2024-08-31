package main

import (
	"io"

	"github.com/spf13/cobra"
)

var refreshDesc = ""

type refreshOptions struct {
	name string
	all  bool
}

func newRefreshCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	// o := &refreshOptions{}
	cmd := &cobra.Command{
		Use:   "refresh NAME",
		Short: "Refresh info for single or all ervers",
		Long:  refreshDesc,
		Run: func(cmd *cobra.Command, args []string) {

			// cmd.Flags().Get
		},
	}

	return cmd
}
