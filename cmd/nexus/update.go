package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

const updateDesc = `
Update APIs managed by Nexus.
`

type updateOptions struct {
	Name string
}

// var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_-]")
// var validUrl = regexp.MustCompile(`https?://[^\s/$.?#].[^\s]*`)

func newUpdateCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	o := &updateOptions{}
	updatedServer := &Server{}

	cmd := &cobra.Command{
		Use:   "update NAME [flags]",
		Short: "update information for a Nexus managed API",
		Long:  updateDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1); err != nil {
				if len(args) == 0 {
					return errors.New("you must specify the name of the API you are updating")
				}

				if len(args) != 1 {
					return errors.New("this command does not take any more arguments")
				}
			}

			o.Name = args[0]

			if !alphanumeric.MatchString(o.Name) {
				return errors.New("name must only contain alphanumeric characters")
			}

			// if !validUrl.MatchString(s.Url) {
			// 	return errors.New("url must be valid url with http or https scheme")
			// }
			return nil

		},
		// ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// 	var comps []string
		// 	if len(args) == 0 {
		// 		comps = cobra.AppendActiveHelp(comps, "You must choose a name for the model serving api you are adding")
		// 	} else if len(args) == 1 {
		// 		comps = cobra.AppendActiveHelp(comps, "You must specify the URL for the model serving api you are adding")
		// 	} else {
		// 		comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments")
		// 	}
		// 	return comps, cobra.ShellCompDirectiveNoFileComp
		// },
		Run: func(cmd *cobra.Command, args []string) {
			for _, server := range *servers {
				if server.Name == o.Name {
					updatedServer = &server

					if cmd.Flag("name").Changed {
						updatedServer.Name, _ = cmd.Flags().GetString("name")

						if !alphanumeric.MatchString(updatedServer.Name) {
							fmt.Fprintf(out, "name must only contain alphanumeric characters\n")
							return
						}
					}

					if cmd.Flag("url").Changed {
						updatedServer.Url, _ = cmd.Flags().GetString("url")

						if !validUrl.MatchString(updatedServer.Url) {
							fmt.Fprintf(out, "url must be valid url with http or https scheme\n")
						}
					}

					if cmd.Flag("api-key").Changed {
						updatedServer.ApiKey, _ = cmd.Flags().GetString("api-key")
					}

					if cmd.Flag("insecure-skip-tls-verify").Changed {
						updatedServer.InsecureSkipTLSVerify, _ = cmd.Flags().GetBool("insecure-skip-tls-verify")
					}

					updatedServer.UpdatedAt = time.Now()

					if err := servers.Update(o.Name, *updatedServer); err != nil {
						fmt.Fprintf(out, "Could not update %s\nnexus: %v\n", o.Name, err.Error())
						return
					}

					fmt.Fprintf(out, "Updated %s\n", o.Name)
					return
				}
			}

			fmt.Fprintf(out, "API %s not found\n", o.Name)
		},
	}

	f := cmd.Flags()

	f.StringP("name", "n", updatedServer.Name, "Updated name for API being updated")
	f.StringP("url", "u", updatedServer.Url, "Updated URL for API being updated")
	f.StringP("api-key", "a", updatedServer.ApiKey, "API Key for connecting to model serving api")
	f.BoolP("insecure-skip-tls-verify", "k", updatedServer.InsecureSkipTLSVerify, "False to skip tls verification; default is true")
	cmd.MarkFlagsOneRequired("name", "url", "api-key", "insecure-skip-tls-verify")

	return cmd
}
