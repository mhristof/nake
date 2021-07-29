package precommit

import (
	"embed"
	"fmt"
	"strings"

	pc "github.com/mhristof/go-precommit"
	"gopkg.in/yaml.v2"
)

//go:embed repos/*
var languages embed.FS

// Repos A list of precommit repositories.
type Repos struct {
	Repos []pc.Repo
}

// Get Return all the precommit repositories rules.
func Get(lang string) []pc.Repo {
	var repos Repos

	config, err := languages.ReadFile(fmt.Sprintf("repos/%s.yaml", strings.ToLower(lang)))
	if err != nil {
		return []pc.Repo{}
	}

	if err := yaml.Unmarshal(config, &repos); err != nil {
		panic(err)
	}

	return repos.Repos
}
