package terraform

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func dirWithFiles(t *testing.T, files map[string]string) (string, func()) {
	dir, err := ioutil.TempDir("", "terraform")

	if err != nil {
		t.Fatal(err)
	}

	for file, content := range files {
		data := []byte(content)
		err := ioutil.WriteFile(filepath.Join(dir, file), data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	return dir, func() {
		_ = os.RemoveAll(dir)
	}
}

func TestAvailable(t *testing.T) {
	var cases = []struct {
		name  string
		files map[string]string
		exp   bool
	}{
		{
			name: "folder without a .tf file",
			exp:  false,
		},
		{
			name: "folder with a .tf file",
			files: map[string]string{
				"main.tf": "",
			},
			exp: true,
		},
	}

	for _, test := range cases {
		dir, cleanup := dirWithFiles(t, test.files)
		defer cleanup()

		var tf = Terraform{
			Pwd: dir,
		}

		assert.Equal(t, test.exp, tf.Available(), test.name)
	}
}
