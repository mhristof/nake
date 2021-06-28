package repo

import (
	"os"
	"path/filepath"

	"github.com/go-enry/go-enry/v2"
)

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

			lang, _ := enry.GetLanguageByExtension(path)
			if _, ok := ignored[lang]; ok {
				return nil
			}

			if _, ok := langs[lang]; !ok {
				ret = append(ret, lang)
			}

			langs[lang] = 1

			return nil
		})
	if err != nil {
		panic(err)
	}

	return ret
}
