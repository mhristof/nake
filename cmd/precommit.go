package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/precommit"
	"github.com/mhristof/nake/repo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	precommitCmd = &cobra.Command{
		Use:   "precommit",
		Short: "Generate pre-commit configuration",
		Run: func(cmd *cobra.Command, args []string) {
			var repos = precommit.Repos{
				Repos: precommit.Get("default"),
			}

			for _, language := range repo.Languages("./") {
				repos.Repos = append(repos.Repos, precommit.Get(language)...)
			}

			reposYAML, err := yaml.Marshal(repos)

			output, err := cmd.Flags().GetString("output")
			if err != nil {
				panic(err)
			}

			log.WithFields(log.Fields{
				"output": output,
			}).Debug("Writing to file")

			if dry {
				fmt.Println(string(reposYAML))
				return
			}

			err = ioutil.WriteFile(output, reposYAML, 0644)
			if err != nil {
				panic(err)
			}

		},
	}
)

func init() {
	precommitCmd.PersistentFlags().StringP("output", "o", ".pre-commit-config.yaml", "Output file to write")

	rootCmd.AddCommand(precommitCmd)
}
