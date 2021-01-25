package cmd

import (
	"github.com/mhristof/nake/bash"
	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/terraform"
	"github.com/spf13/cobra"
)

var (
	planCmd = &cobra.Command{
		Use:   "plan",
		Short: "Terraform plan",
		Run: func(cmd *cobra.Command, args []string) {
			var tf = terraform.Terraform{
				Pwd: pwd,
			}

			if !tf.Available() {
				log.Fatal("Terraform not available for this folder")
			}

			commands := []string{
				tf.Init(),
				tf.Plan(),
			}
			for _, command := range commands {
				if command == "" {
					continue
				}
				bash.Eval(command)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(planCmd)
}
