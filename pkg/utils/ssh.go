package utils

"os/exec"

func SSHCommand(IP string, cmd ...string) *exec.Cmd {
	arg := append(
		[]string{"-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile /dev/null", "-i", filepath.Join(os.Getwd(), "ssh/id_rsa"),
			fmt.Sprintf("root@%s", IP),
		},
		cmd...,
	)
	return exec.Command("ssh", arg...)
}