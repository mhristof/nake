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

func modTime(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		panic(err)
	}

	return info.ModTime()
}

func (t *Terraform) Plan() string {
	plan := filepath.Join(t.Pwd, "terraform.tfplan")
	planInfo, err := os.Stat(plan)
	// if the file doesnt exist, use beginning of Epoch as the plan timestamp
	// to force a rebuild
	planMod := time.Time{}
	if err == nil {
		planMod = planInfo.ModTime()
	}

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
