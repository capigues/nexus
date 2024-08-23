package main

import (
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/spf13/cobra"
)

const addDesc = `
ADD MORE
`

var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_-]*$")
var validUrl = regexp.MustCompile(`https?://[^\s/$.?#].[^\s]*`)

type addOptions struct {
	name string
	url  string
}

func newAddCommand(servers *ModelServers, out io.Writer) *cobra.Command {
	o := &addOptions{}

	cmd := &cobra.Command{
		Use:   "add NAME URL",
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

			name, url := args[0], args[1]

			if !alphanumeric.MatchString(name) {
				return errors.New("name must only contain alphanumeric characters")
			}

			if !validUrl.MatchString(url) {
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
			o.name = args[0]
			o.url = args[1]

			if err := servers.Add(o.name, o.url); err != nil {
				fmt.Fprintf(out, "Could not add %s\nnexus: %v", o.name, err.Error())
				return
			}

			fmt.Fprintf(out, "Added %s\n", o.name)
		},
	}

	return cmd
}
