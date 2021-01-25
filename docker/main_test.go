package docker

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/nake/bash"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	var cases = []struct {
		name  string
		files map[string]string
		exp   string
	}{
		{
			name: "Simple dockerfile",
			files: map[string]string{
				"Dockerfile": heredoc.Doc(`FROM golang`),
			},
			exp: "docker build --tag {{.Image}} --file Dockerfile .",
		},
	}

	for _, test := range cases {
		dir, cleanup := bash.DirWithFiles(t, test.files)
		defer cleanup()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}

		var d = Docker{
			Dockerfile: "Dockerfile",
			Pwd:        ".",
		}

		image := filepath.Base(dir)
		buildTemplate, err := template.New("build").Parse(test.exp)
		if err != nil {
			t.Fatal(err)
		}

		var build bytes.Buffer
		err = buildTemplate.Execute(&build, struct{ Image string }{Image: image})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, build.String(), d.Build(image), test.name)
	}
}
