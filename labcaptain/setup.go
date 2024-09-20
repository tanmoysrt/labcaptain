package main

import (
	_ "embed"
	"os"
	"os/exec"
)

//go:embed scripts/local_setup.sh
var localSetupScript string

func setup() error {
	// Run the script in the server
	cmd := exec.Command("/bin/bash", "-c", localSetupScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
