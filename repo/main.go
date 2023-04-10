package repo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-enry/go-enry/v2"
	log "github.com/sirupsen/logrus"
)

// Languages Return a list of all languages filetypes inside `dest` folder.
func Languages(dest string, ignore []string) []string {
	var ret []string

	ignored := map[string]bool{
		"":            true,
		"Ignore List": true,
	}
	langs := make(map[string]int)

	ignore = append(ignore, ".terraform")
	ignore = append(ignore, ".terragrunt-cache")
	ignore = append(ignore, "env.hcl")

	err := filepath.Walk(dest,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			for _, dir := range ignore {
				if strings.Contains(path, dir) {
					return nil
				}
			}

			if strings.Contains(path, ".git") {
				langs["git"] = 1

				return nil
			}

			if strings.HasSuffix(path, ".pre-commit-config.yaml") {
				langs["precommit"] = 1

				return nil
			}

			if strings.HasSuffix(path, "terragrunt.hcl") {
				langs["terragrunt"] = 1
				log.WithFields(log.Fields{
					"path": path,
				}).Debug("found terragrunt file")

				return nil
			}

			if strings.HasPrefix(filepath.Base(path), "Dockerfile") {
				langs["Docker"] = 1

				log.WithFields(log.Fields{
					"path": path,
				}).Debug("found Dockerfile")

				return nil
			}

			lang, _ := enry.GetLanguageByExtension(path)
			if _, ok := ignored[lang]; ok {
				return nil
			}

			log.WithFields(log.Fields{
				"lang": lang,
				"path": path,
			}).Debug("found language")

			langs[lang] = 1

			return nil
		})
	if err != nil {
		panic(err)
	}

	for lang := range langs {
		ret = append(ret, lang)
	}

	return ret
}
