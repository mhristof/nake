package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
