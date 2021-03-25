package cmd

import (
	"fmt"

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
			var repos precommit.Repos

			for _, language := range repo.Languages("./") {
				repos.Repos = append(repos.Repos, precommit.Get(language)...)
			}

			reposJSON, err := yaml.Marshal(repos)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(reposJSON))
		},
	}
)

func init() {
	rootCmd.AddCommand(precommitCmd)
}
