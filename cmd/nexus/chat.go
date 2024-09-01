package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const chatDesc = `
Chat with model hosted by an API managed by Nexus
`

type chatOptions struct {
	Name        string
	Temperature int
	MaxTokens   int
	// Stream? bool
}

func newChatCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	o := &chatOptions{}

	cmd := &cobra.Command{
		Use:   "chat NAME",
		Short: "chat with an api managed by Nexus",
		Long:  chatDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("you must specify an API to chat with")
			}

			if len(args) != 1 {
				return errors.New("can only specify one API to chat with")
			}

			o.Name = args[0]

			if !servers.Find(o.Name) {
				return fmt.Errorf("API %v not found", o.Name)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			o.Temperature = *cmd.Flags().IntP("temperature", "t", 0, "temperature for LLM chat")
			o.MaxTokens = *cmd.Flags().IntP("max-tokens", "m", 1024, "max tokens that the LLM will generate")

			for _, server := range *servers {
				if server.Name == o.Name {
					err := server.Chat(o.Temperature, o.MaxTokens)
					if err != nil {
						fmt.Fprintf(out, "Could not start chat with %s\nnexus: %v\n", server.Name, err.Error())
					}

					return
				}
			}

			fmt.Fprintf(out, "API %v not found\n", o.Name)
		},
	}

	return cmd
}
