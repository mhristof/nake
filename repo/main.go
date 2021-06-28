package repo

import (
	"io/ioutil"

	"github.com/go-enry/go-enry/v2"
)

func Languages(dest string) []string {
	var ret []string

	ignored := map[string]bool{
		"":            true,
		"Ignore List": true,
	}
	langs := make(map[string]int)

	files, err := ioutil.ReadDir(dest)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		lang, _ := enry.GetLanguageByExtension(file.Name())
		if _, ok := ignored[lang]; ok {
			continue
		}

		if _, ok := langs[lang]; !ok {
			ret = append(ret, lang)
		}

		langs[lang] = 1
	}

	return ret
}
