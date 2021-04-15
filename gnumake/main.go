package gnumake

import (
	"fmt"

	"github.com/mhristof/nake/log"
	"github.com/mhristof/nake/repo"
)

type Rules []Rule
type Rule struct {
	Targets       string
	Prerequisites string
	Recipe        string
	Phony         bool
	Help          string
}

var RulesLib = map[string][]Rule{
	"Python": []Rule{
		Rule{
			Help:    "Run pep8 for the current directory",
			Targets: "pep8",
			Recipe:  "pycodestyle --ignore=E501",
			Phony:   true,
		},
	},
	"Terraform": []Rule{
		Rule{
			Help:          "Force run 'terraform init'",
			Targets:       "init",
			Prerequisites: ".terraform",
			Phony:         true,
		},
		Rule{
			Targets: ".terrform",
			Recipe:  "terraform init",
		},
		Rule{
			Help:          "Runs 'terraform plan'",
			Targets:       "plan",
			Prerequisites: "terraform.tfplan",
			Phony:         true,
		},
		Rule{
			Help:          "Creates terraform.tfplan if required",
			Targets:       "terraform.tfplan",
			Prerequisites: "$(shell find ./ -name '*.tf')",
			Recipe:        "terraform plan -out $@",
		},
		Rule{
			Help:          "Run 'terraform apply'",
			Targets:       "apply",
			Prerequisites: "terraform.tfplan",
			Phony:         true,
		},
		Rule{
			Help:          "Run 'terraform apply' if required'",
			Targets:       "terraform.tfstate",
			Prerequisites: "terraform.tfplan",
			Recipe:        "terraform apply terraform.tfplan",
		},
		Rule{
			Help:    "Forcefully update terraform state",
			Targets: "force",
			Recipe:  "touch *.tf && make terraform.tfstate",
			Phony:   true,
		},
		Rule{
			Help:    "Run 'terraform destroy'",
			Targets: "destroy",
			Recipe:  "terraform destroy -auto-approve && rm -f terraform.tf{state,plan}",
			Phony:   true,
		},
		Rule{
			Help:          "Clean the repository resources",
			Targets:       "clean",
			Prerequisites: "destroy",
			Phony:         true,
		},
	},
}

func Generate() {
	var rules Rules

	for _, language := range repo.Languages("./") {
		log.WithFields(log.Fields{
			"language": language,
		}).Debug("Found")

		rules = append(rules, RulesLib[language]...)
	}

	rules.Print()
}

func (r Rules) Print() {
	for _, rule := range r {
		if rule.Phony {
			fmt.Println(fmt.Sprintf(".PHONY: %s", rule.Targets))
		}
		fmt.Println(fmt.Sprintf("%s: %s ## %s", rule.Targets, rule.Prerequisites, rule.Help))

		if rule.Recipe != "" {
			fmt.Println(fmt.Sprintf("\t%s", rule.Recipe))
		}
		fmt.Println("")
	}
}
