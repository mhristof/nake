package cmd

import (
	"github.com/mhristof/nake/bash"
	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/terraform"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Terraform init",
		Run: func(cmd *cobra.Command, args []string) {
			var tf = terraform.Terraform{
				Pwd: pwd,
			}

			if !tf.Available() {
				log.Fatal("Terraform not available for this folder")
			}

			command := tf.Init()

			if command != "" {
				bash.Eval(command)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}
