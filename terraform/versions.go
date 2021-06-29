package terraform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/nake/log"
	"golang.org/x/mod/semver"
)

type versionTF struct {
	RequiredVersion   string              `json:"required_version,omitempty"`
	RequiredProviders map[string]provider `json:"required_providers,omitempty"`
}

type provider struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

var versionsTFTemplate = heredoc.Doc(`
	terraform {
	  {{ if .RequiredVersion -}} required_version = "{{ .RequiredVersion }}"{{ end }}
	  required_providers {
	    {{ range $key, $val := .RequiredProviders -}}
	    {{ $key }} = {
	      source = "{{ $val.Source }}"
	      version = "{{ $val.Version }}"
	    }
	    {{end}}
	  }
	}
`)

type binaryVersion struct {
	Platform           string      `json:"platform"`
	ProviderSelections interface{} `json:"provider_selections"`
	TerraformOutdated  bool        `json:"terraform_outdated"`
	TerraformVersion   string      `json:"terraform_version"`
}

func terraformVersion() string {
	var stdout bytes.Buffer

	cmd := exec.Command("terraform", "version", "-json")
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	var tfv binaryVersion

	err = json.Unmarshal(stdout.Bytes(), &tfv)
	if err != nil {
		panic(err)
	}

	return tfv.TerraformVersion
}

// Versions Generate versions.tf content for the given `source` directory
func Versions(source string, strict bool) string {
	var ver = versionTF{
		RequiredVersion:   terraformVersion(),
		RequiredProviders: make(map[string]provider),
	}

	if !strict {
		ver.RequiredVersion = "~> " + ver.RequiredVersion
	}

	providers := fmt.Sprintf("%s/providers/registry.terraform.io/hashicorp", source)

	if info, err := os.Stat(providers); err == nil && !info.IsDir() {
		log.WithFields(log.Fields{
			"providers": providers,
		}).Error("Not a dir")

		return ""
	}

	err := filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if path == source {
				return nil
			}

			fields := strings.Split(path, "/")
			version := fields[len(fields)-1]
			if !semver.IsValid("v" + version) {
				return nil
			}

			if !strict {
				version = "~> " + version
			}

			ver.RequiredProviders[fields[len(fields)-2]] = provider{
				Source:  filepath.Join(fields[len(fields)-3 : len(fields)-1]...),
				Version: version,
			}

			return nil
		})

	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("versions.tf").Parse(versionsTFTemplate))

	b := new(strings.Builder)

	err = t.Execute(b, ver)
	if err != nil {
		panic(err)
	}

	return b.String()
}
