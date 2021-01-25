package bash

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Eval(command string) error {
	cmd := exec.Command("/bin/bash", "-c", command)

	fmt.Println(fmt.Sprintf("cmd.Args: %+v", cmd.Args))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func DirWithFiles(t *testing.T, files map[string]string) (string, func()) {
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
