package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const listDesc = `
ADD MORE
`

func newListCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all saved apis",
		Long:  listDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}

			if err := cobra.ExactArgs(0); err != nil {
				return errors.New("this command does not take any more arguments")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := servers.List(); err != nil {
				fmt.Fprintf(out, "Could not list apis\nnexus: %v", err.Error())
				return
			}
		},
	}

	return cmd
}
