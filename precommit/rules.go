package precommit

var languages = map[string][]Repo{
	"Go": []Repo{
		{
			Repo: "https://github.com/tekwizely/pre-commit-golang",
			Rev:  "master",
			Hooks: []Hook{
				Hook{ID: "go-build-mod"},
				Hook{ID: "go-build-pkg"},
				Hook{ID: "go-build-repo-mod"},
				Hook{ID: "go-build-repo-pkg"},
				Hook{ID: "go-test-mod"},
				Hook{ID: "go-test-pkg"},
				Hook{ID: "go-test-repo-mod"},
				Hook{ID: "go-test-repo-pkg"},
				Hook{ID: "go-vet"},
				Hook{ID: "go-vet-mod"},
				Hook{ID: "go-vet-pkg"},
				Hook{ID: "go-vet-repo-mod"},
				Hook{ID: "go-vet-repo-pkg"},
				Hook{ID: "go-sec-mod"},
				Hook{ID: "go-sec-pkg"},
				Hook{ID: "go-sec-repo-mod"},
				Hook{ID: "go-sec-repo-pkg"},
				Hook{ID: "go-fmt"},
				Hook{ID: "go-imports"},
				Hook{ID: "go-returns"},
				Hook{ID: "go-lint"},
				Hook{ID: "go-critic"},
				Hook{ID: "golangci-lint"},
				Hook{ID: "golangci-lint-mod"},
				Hook{ID: "golangci-lint-pkg"},
				Hook{ID: "golangci-lint-repo-mod"},
				Hook{ID: "golangci-lint-repo-pkg"},
			},
		},
	},
}

func Get(lang string) []Repo {
	value, ok := languages[lang]
	if ok != true {
		return []Repo{}
	}
	return value
}
