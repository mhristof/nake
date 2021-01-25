package cmd

import (
	"github.com/mhristof/nake/bash"
	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/terraform"
	"github.com/spf13/cobra"
)

var (
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Terraform apply",
		Run: func(cmd *cobra.Command, args []string) {
			var tf = terraform.Terraform{
				Pwd: pwd,
			}

			if !tf.Available() {
				log.Fatal("Terraform not available for this folder")
			}

			force, _ := cmd.Flags().GetBool("force")
			commands := []string{
				tf.Init(),
				tf.Plan(),
				tf.Apply(force),
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
	applyCmd.PersistentFlags().BoolP("force", "f", false, "Force 'terraform apply' ignoring source files")
	rootCmd.AddCommand(applyCmd)
}
