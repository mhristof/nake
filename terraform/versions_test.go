package terraform

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createfs(t *testing.T, files []string) (string, func()) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		path := filepath.Join(dir, file)
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

	}
	return dir, func() {
		os.RemoveAll(dir)
	}
}

func TestVersions(t *testing.T) {
	var cases = []struct {
		name   string
		strict bool
		fs     []string
	}{
		{
			name: "strict off",
			fs: []string{
				".terraform/providers/registry.terraform.io/hashicorp/time/0.7.1",
				".terraform/providers/registry.terraform.io/hashicorp/local/2.1.0",
			},
		},
		{
			name:   "strict",
			strict: true,
			fs: []string{
				".terraform/providers/registry.terraform.io/hashicorp/time/0.7.1",
				".terraform/providers/registry.terraform.io/hashicorp/local/2.1.0",
			},
		},
	}

	for _, test := range cases {
		source, cleanup := createfs(t, test.fs)
		out, err := ioutil.ReadFile(filepath.Join("fixtures", strings.ReplaceAll(test.name, " ", ".")+".tf"))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, string(out), Versions(source, test.strict), test.name)
		cleanup()
	}
}
