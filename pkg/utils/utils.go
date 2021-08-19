package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	//	"time"
)

func ChangeXMLSpec(path string) error {
	fmt.Println(path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Println(data)
	return nil
}

func CopyRawDisks(remotehostIP string, ftpservIP string, distro string) (string, error) {
	_, xml := filepath.Split(fmt.Sprintf("http://%s/sle%s-fake-NEW.xml", ftpservIP, distro))
	//fmt.Printf("FULL PATH: %s\n", fmt.Sprintf("http://%s/sle%s-fake-NEW.xml", ftpservIP, distro))
	//time.Sleep(30 * time.Second)
	command := []string{"wget", fmt.Sprintf("http://%s/sle%s_fake_baremetal_xenvirthost_client.qcow2", ftpservIP, distro), fmt.Sprintf("http://%s/%s", remotehostIP, xml), "-P", "/var/lib/libvirt/images/"}
	//fmt.Println(command)
	//time.Sleep(30 * time.Second)
	cmd := SSHCommand(remotehostIP, command...)
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	NiceBuffRunner(cmd, pwd)
	return filepath.Join("/var/lib/libvirt/images/", xml), nil
}

func (cluster *Command) SaveJSN(rootdir string) error {
	file, err := json.MarshalIndent(cluster, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(rootdir, "cluster.json"), file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func OpenJSN(filelocation string) (cluster *Command, err error) {
	file, err := os.Open(filelocation)
	defer file.Close()
	if err != nil {
		return
	}
	if err = json.NewDecoder(file).Decode(&cluster); err != nil {
		return
	}
	//fmt.Printf("CLUSTER: %v\n", cluster)
	return
}

func SliceExec(command []string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)
	return cmd
}

func NiceBuffRunner(cmd *exec.Cmd, workdir string) (string, string) {
	var stdoutBuf, stderrBuf bytes.Buffer
	//newEnv := append(os.Environ(), ENV...)
	//cmd.Env = newEnv
	cmd.Dir = workdir
	pipe, _ := cmd.StdoutPipe()
	errpipe, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		return fmt.Sprintf("%s", os.Stdout), fmt.Sprintf("%s", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, pipe)
		wg.Done()
	}()
	go func() {
		_, errStderr = io.Copy(stderr, errpipe)
		wg.Wait()
	}()
	err = cmd.Wait()
	if err != nil {
		return fmt.Sprintf("%s", os.Stdout), fmt.Sprintf("%s", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("Command runninng error: failed to capture stdout or stderr\n")
	}
	return stdoutBuf.String(), stderrBuf.String()
}
