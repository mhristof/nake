package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mhristof/nake/log"
)

type Terraform struct {
	Pwd string
}

func (t *Terraform) Available() bool {
	files, err := ioutil.ReadDir(t.Pwd)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".tf") {
			return true
		}
	}
	return false
}

func (t *Terraform) Init() string {
	dotTerraform := filepath.Join(t.Pwd, ".terraform")

	if _, err := os.Stat(dotTerraform); !os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"dotTerraform": dotTerraform,
		}).Debug(".terraform already exists")
		return ""
	}

	return fmt.Sprintf("cd %s && terraform init", t.Pwd)
}

// modTime Return the modification time of a file. If it doesnt exist, report
// Epoch beginning of time.
func modTime(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		return time.Time{}
	}

	return info.ModTime()
}

func (t *Terraform) Plan() string {
	planMod := modTime(filepath.Join(t.Pwd, "terraform.tfplan"))

	files, err := ioutil.ReadDir(t.Pwd)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		path := filepath.Join(t.Pwd, file.Name())
		fileMod := modTime(path)
		if fileMod.After(planMod) {
			return "terraform plan"
		}
	}

	return ""
}

func (t *Terraform) Apply(force bool) string {
	stateMod := modTime(filepath.Join(t.Pwd, "terraform.tfstate"))
	planMod := modTime(filepath.Join(t.Pwd, "terraform.tfplan"))

	if force || planMod.After(stateMod) {
		return "terraform apply terraform.tfplan"
	}

	return ""
}
