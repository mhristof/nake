package cmd

import (
	"github.com/mhristof/nake/bash"
	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/terraform"
	"github.com/spf13/cobra"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clea",
		Short: "Terraform clean",
		Run: func(cmd *cobra.Command, args []string) {
			var tf = terraform.Terraform{
				Pwd: pwd,
			}

			if !tf.Available() {
				log.Fatal("Terraform not available for this folder")
			}

			command := tf.Clean()

			if command != "" {
				bash.Eval(command)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}
