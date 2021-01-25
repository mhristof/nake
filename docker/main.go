package docker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mhristof/nake/log"
)

type Docker struct {
	Dockerfile string
	Pwd        string
}

func image(dockerfile string) string {
	return filepath.Base(filepath.Dir(dockerfile))
}

func (d *Docker) Build(image string) string {
	if _, err := os.Stat(d.Dockerfile); os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"d.Dockerfile": d.Dockerfile,
		}).Panic("Dockerfile missing")
	}

	return fmt.Sprintf("docker build --tag %s --dockerfile %s %s", image, d.Dockerfile, d.Pwd)
}
