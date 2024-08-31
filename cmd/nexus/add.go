package main

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/spf13/cobra"
)

const addDesc = `
Add API to be managed by Nexus
`

var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_-]*$")
var validUrl = regexp.MustCompile(`https?://[^\s/$.?#].[^\s]*`)

func newAddCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	s := &Server{}

	cmd := &cobra.Command{
		Use:   "add NAME URL [flags]",
		Short: "add a model serving api at url",
		Long:  addDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(2); err != nil {
				if len(args) == 0 {
					return errors.New("you must choose a name for the model serving api you are adding")
				}
				if len(args) == 1 {
					return errors.New("you must specify the URL for the model serving api you are adding")
				}

				if len(args) != 2 {
					return errors.New("this command does not take any more arguments")
				}
			}

			s.Name = args[0]
			s.Url = args[1]

			if !alphanumeric.MatchString(s.Name) {
				return errors.New("name must only contain alphanumeric characters")
			}

			if !validUrl.MatchString(s.Url) {
				return errors.New("url must be valid url with http or https scheme")
			}
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
			s.Name = args[0]
			s.Url = args[1]

			s.ApiKey, _ = cmd.Flags().GetString("api-key")
			s.InsecureSkipTLSVerify, _ = cmd.Flags().GetBool("insecure-skip-tls-verify")

			s.UpdatedAt = time.Now()
			s.CreatedAt = time.Now()

			if err := servers.Add(s); err != nil {
				fmt.Fprintf(out, "Could not add %s\nnexus: %v\n", s.Name, err)
				return
			}

			fmt.Fprintf(out, "Added %s\n", s.Name)
		},
	}

	f := cmd.Flags()

	f.StringP("api-key", "a", "", "API Key for connecting to model serving api")
	f.BoolP("insecure-skip-tls-verify", "k", false, "False to skip tls verification; default is true")

	return cmd
}
