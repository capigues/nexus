package main

import (
	"io"

	"github.com/spf13/cobra"
)

const rootDesc = `
Nexus is a API federation tool built for managing, monitoring and querying single model serving APIs from one place. Currently only support OpenAI single model serving runtimes.
`

func NewRootCmd(servers *ModelServers, out io.Writer, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "nexus",
		Short: "model serving api management tool",
		Long:  rootDesc,
	}

	cmd.AddCommand(
		newAddCommand(servers, out), // Get more info from api route for servers.json (model name, size, that type of stuff); will be update with each refresh
		newRemoveCommand(servers, out),
		newListCommand(servers, out),    // Need to fix this. Find better place to check connectivity (calls will take to long if everytime we call)
		newRefreshCommand(servers, out), // Refresh command to sync url with their api and use that to update statuses; also get more info to add to servers.json
		newUpdateCommand(servers, out),
		// newChatCommand(servers, out), // Start quick chat with api selected
		// newServeCommand(servers, out), // Serve single endpoint to chat across multiple openAI API endpoints
		// newInfoCommand(servers, out), // Specific api information (idk)
	)

	return cmd, nil
}
