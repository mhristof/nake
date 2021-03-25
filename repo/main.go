package repo

import (
	"io/ioutil"

	"github.com/go-enry/go-enry/v2"
)

func Languages(dest string) []string {
	ignored := map[string]bool{
		"":            true,
		"Ignore List": true,
	}
	langs := make(map[string]int)
	var ret []string

	files, err := ioutil.ReadDir(dest)
	if err != nil {
		panic(err)
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		lang, _ := enry.GetLanguageByExtension(file.Name())
		if _, ok := ignored[lang]; ok == true {
			continue
		}

		if _, ok := langs[lang]; ok != true {
			ret = append(ret, lang)
		}
		langs[lang] = 1
	}

	return ret
}
