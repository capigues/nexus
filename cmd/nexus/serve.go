package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const serveDesc = `
Serve all APIs saved in Nexus at one single endpoint
`

type serveOptions struct {
	Port string
}

func newServeCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	o := serveOptions{}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "start nexus proxy server",
		Long:  serveDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not take any arguments")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Port = *cmd.Flags().StringP("port", "p", "8080", "port for Nexus server")

			err := servers.Serve(out, o.Port)
			if err != nil {
				fmt.Fprintf(out, "%v\n", err.Error())
			}
		},
	}

	return cmd
}
