package cmd

import (
	"io/ioutil"
	"path/filepath"

	"github.com/mhristof/nake/terraform"
	"github.com/spf13/cobra"
)

var (
	versionsTFCmd = &cobra.Command{
		Use:   "versions.tf",
		Short: "Generate versions.tf file",
		Run: func(cmd *cobra.Command, args []string) {
			// do something
			source, err := cmd.Flags().GetString("source")
			if err != nil {
				panic(err)
			}

			strict, err := cmd.Flags().GetBool("strict")
			if err != nil {
				panic(err)
			}

			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				panic(err)
			}

			output, err := cmd.Flags().GetString("output")
			if err != nil {
				panic(err)
			}

			output = filepath.Join(dir, output)

			err = ioutil.WriteFile(output, []byte(terraform.Versions(source, strict)), 0644)
			if err != nil {
				panic(err)
			}

		},
	}
)

func init() {
	versionsTFCmd.PersistentFlags().StringP("source", "s", ".terraform", "The .terraform directory to scan for versions")
	versionsTFCmd.PersistentFlags().BoolP("strict", "x", false, "Strict mode, wont add ~> in the version found.")
	versionsTFCmd.PersistentFlags().StringP("output", "o", "version.tf", "Output file")

	rootCmd.AddCommand(versionsTFCmd)
}
