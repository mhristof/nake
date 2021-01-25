package bash

import (
	"fmt"
	"os"
	"os/exec"
)

func Eval(command string) error {
	cmd := exec.Command("/bin/bash", "-c", command)

	fmt.Println(fmt.Sprintf("cmd.Args: %+v", cmd.Args))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
