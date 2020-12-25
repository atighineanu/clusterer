package utils

import (
	"os/exec"
	"os"
	"path/filepath"
	"fmt"
)


func SSHCommand(IP string, cmd ...string) *exec.Cmd {
	workdir, err := os.Getwd()
	if err != nil {
		fmt.Printf("pkg/utils/ssh.go ERROR: %v\n", err)
	}
	arg := append(
		[]string{"-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile /dev/null", "-i", filepath.Join(workdir, "ssh/id_rsa"),
			fmt.Sprintf("root@%s", IP),
		},
		cmd...,
	)
	return exec.Command("ssh", arg...)
}