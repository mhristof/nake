package precommit

import (
	"embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed repos/*
var languages embed.FS

func Get(lang string) []Repo {
	var repos Repos

	config, err := languages.ReadFile(fmt.Sprintf("repos/%s.yaml", strings.ToLower(lang)))
	if err != nil {
		return []Repo{}
	}

	if err := yaml.Unmarshal(config, &repos); err != nil {
		panic(err)
	}

	return repos.Repos
}
