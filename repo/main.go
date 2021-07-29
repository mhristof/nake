package repo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-enry/go-enry/v2"
)

// Languages Return a list of all languages filetypes inside `dest` folder.
func Languages(dest string) []string {
	var ret []string

	ignored := map[string]bool{
		"":            true,
		"Ignore List": true,
	}
	langs := make(map[string]int)

	err := filepath.Walk(dest,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if strings.Index(path, ".git") >= 0 {
				langs["git"] = 1
				return nil
			}

			if strings.Index(path, ".terraform") > 0 {
				return nil
			}

			if strings.HasSuffix(path, ".pre-commit-config.yaml") {
				langs["precommit"] = 1

				return nil
			}

			lang, _ := enry.GetLanguageByExtension(path)
			if _, ok := ignored[lang]; ok {
				return nil
			}

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
