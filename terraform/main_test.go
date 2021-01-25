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

func TestPlan(t *testing.T) {
	var cases = []struct {
		name  string
		files map[string]string
		exp   string
	}{
		{
			name: "missing terraform.tfplan",
			files: map[string]string{
				"main.tf": "",
			},
			exp: "terraform plan",
		},
		{
			name: "up2date plan file",
			files: map[string]string{
				"main.tf": "",
				// terraform.tfplan is generatred after main.tf
				"terraform.tfplan": "",
			},
			exp: "",
		},
		{
			name: "old plan file",
			files: map[string]string{
				"terraform.tfplan": "",
				// main.tf is newer than plan file.
				"main.tf": "",
			},
			exp: "terraform plan",
		},
	}

	for _, test := range cases {
		dir, cleanup := dirWithFiles(t, test.files)
		defer cleanup()

		var tf = Terraform{
			Pwd: dir,
		}

		assert.Equal(t, test.exp, tf.Plan(), test.name)
	}
}

func TestApply(t *testing.T) {
	var cases = []struct {
		name  string
		files map[string]string
		force bool
		exp   string
	}{
		{
			name: "Plan file newer than statefile",
			files: map[string]string{
				"main.tf":           "",
				"terraform.tfstate": "",
				"terraform.tfplan":  "",
			},
			exp: "terraform apply terraform.tfplan && rm terraform.tfplan",
		},
		{
			name: "Plan file older than statefile",
			files: map[string]string{
				"main.tf":           "",
				"terraform.tfplan":  "",
				"terraform.tfstate": "",
			},
			exp: "",
		},
		{
			name: "Force apply",
			files: map[string]string{
				"main.tf":           "",
				"terraform.tfplan":  "",
				"terraform.tfstate": "",
			},
			force: true,
			exp:   "terraform apply terraform.tfplan && rm terraform.tfplan",
		},
	}

	for _, test := range cases {
		dir, cleanup := dirWithFiles(t, test.files)
		defer cleanup()

		var tf = Terraform{
			Pwd: dir,
		}

		assert.Equal(t, test.exp, tf.Apply(test.force), test.name)
	}
}

func TestInit(t *testing.T) {
	var cases = []struct {
		name  string
		files map[string]string
		exp   string
	}{
		{
			name: "folder without a .terraform file",
			exp:  "terraform init",
		},
		{
			name: "folder with a .terraform file",
			files: map[string]string{
				".terraform": "",
			},
			exp: "",
		},
	}

	for _, test := range cases {
		dir, cleanup := dirWithFiles(t, test.files)
		defer cleanup()

		var tf = Terraform{
			Pwd: dir,
		}

		assert.Equal(t, test.exp, tf.Init(), test.name)
	}
}
