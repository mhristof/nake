package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/precommit"
	"github.com/mhristof/nake/repo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	precommitCmd = &cobra.Command{
		Use:   "precommit",
		Short: "Generate pre-commit configuration",
		Run: func(cmd *cobra.Command, args []string) {
			var repos = precommit.Repos{
				Repos: precommit.Get("default"),
			}

			dir, err := rootCmd.PersistentFlags().GetString("dir")
			if err != nil {
				panic(err)
			}

			ignore, err := rootCmd.PersistentFlags().GetStringSlice("ignore")
			if err != nil {
				panic(err)
			}

			for _, language := range repo.Languages(dir, ignore) {
				log.WithFields(log.Fields{
					"language": language,
				}).Debug("Adding precommit rules")
				repos.Repos = append(repos.Repos, precommit.Get(language)...)
			}

			var b bytes.Buffer
			yamlEncoder := yaml.NewEncoder(&b)
			yamlEncoder.SetIndent(2) // this is what you're looking for
			err = yamlEncoder.Encode(&repos)
			if err != nil {
				panic(err)

			}

			output, err := cmd.Flags().GetString("output")
			if err != nil {
				panic(err)
			}

			outputAbs, err := filepath.Abs(output)
			if err != nil {
				panic(err)
			}

			if outputAbs != output {
				output = filepath.Join(dir, output)
			}

			log.WithFields(log.Fields{
				"output": output,
			}).Debug("Writing to file")

			if dry {
				fmt.Println(b.String())
				return
			}

			err = ioutil.WriteFile(output,
				bytes.Join([][]byte{
					[]byte("---"),
					b.Bytes()},

					[]byte("\n"),
				),
				0644,
			)
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
