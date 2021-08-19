package libvirtd

import (
	"clusterer/pkg/utils"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func DefineVMfromXML(path string, remotehostIP string) (string, error) {
	command := []string{"virsh", "define", path}
	//cmd := exec.Command(command[0], command[1:]...)
	cmd := utils.SSHCommand(remotehostIP, command...)
	out, _ := utils.NiceBuffRunner(cmd, "./")
	tmp := strings.Split(out, "\n")
	var machine string
	for _, line := range tmp {
		if strings.Contains(line, "defined") {
			chunkedline := strings.Split(line, " ")
			if strings.Contains(chunkedline[0], "Domain") && strings.Contains(chunkedline[1], "xenvirthost") {
				machine = chunkedline[1]
				fmt.Printf("Machine: %s\n", machine)
			}
		}
	}
	return machine, nil
}

func CloneVol(cluster utils.Command, seed string, machine string) error {
	log.Println("Clonning Volume(s)...")
	cmdstring := []string{"sudo", "virsh", "vol-clone", seed, machine, "--pool", cluster.Pool.Name}
	command := utils.SliceExec(cmdstring)
	_, err := utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	return nil
}

func CloneVM(cluster utils.Command, seed string, machine string) error {
	log.Println("Clonning VM(s)...")
	cmdstring := []string{"sudo", "virt-clone", "-o", seed, "-n", machine, "--preserve-data", "-f", filepath.Join(cluster.Pool.Path, machine)}
	command := utils.SliceExec(cmdstring)
	_, err := utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	return nil
}

func StartVM(machine string, remote string) error {
	log.Println("Starting VM(s)...")
	cmdstring := []string{"sudo", "virsh", "start", machine}
	if remote == "" {
		command := utils.SliceExec(cmdstring)
		_, err := utils.NiceBuffRunner(command, "/home/user")
		if err != "" {
			return errors.New(err)
		}
	} else {
		cmd := utils.SSHCommand(remote, cmdstring...)
		utils.NiceBuffRunner(cmd, "./")
	}
	return nil
}

func WaitForIP(cluster utils.Command, machine, remote string) error {
	CheckIfExists(cluster, machine, remote, true)
	return nil
}

func CmdRemoteSpitter(cmd []string, remote string) *exec.Cmd {
	if remote == "" {
		return exec.Command(cmd[0], cmd[1:]...)
	} else {
		return utils.SSHCommand(remote, cmd...)
	}
}

func CheckIfExists(cluster utils.Command, machine, remote string, silent bool) (utils.Command, error) {
	inc := 0
	for {
		cmdstring := []string{"sudo", "virsh", "domifaddr", machine, "--source", "agent"}
		comm := CmdRemoteSpitter(cmdstring, remote)
		resp, err := comm.CombinedOutput()
		if !silent {
			fmt.Println(fmt.Sprintf("%s", string(resp)))
		}
		if err != nil {
			cmdstring2 := []string{"sudo", "virsh", "list"}
			comm2 := CmdRemoteSpitter(cmdstring2, remote)
			resp2, err2 := comm2.CombinedOutput()
			if err2 != nil {
				return cluster, err2
			}
			if !silent {
				fmt.Println(fmt.Sprintf("%s", string(resp2)))
			}
			if strings.Contains(fmt.Sprintf("%s", string(resp2)), machine) {
				log.Println("Machine not ready yet...")
				time.Sleep(3 * time.Second)
			} else {
				log.Println("Looks like machine doesn't exist or offline...Will see if we remove it from cluster catalogue...")
				SeeIfOffline(cluster, machine, remote, silent)
			}
		} else {
			found := false
			for _, row := range strings.Split(fmt.Sprintf("%s", string(resp)), "\n") {
				if (strings.Contains(row, "eth") || strings.Contains(row, "br")) && strings.Contains(row, "ipv4") {
					reg := regexp.MustCompile(`\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}/\d{2}`)
					//mac := regexp.MustCompile()
					if len(reg.FindStringSubmatch(row)) > 0 {
						var temp net.HardwareAddr
						for _, value := range strings.Split(row, " ") {
							temp, err = net.ParseMAC(value)
							if err == nil {
								break
							}
						}
						log.Printf("Found a proper IP for the newly spawned VM...\n(libvirt) Machine name: %s\nIP: %s\nMAC: %v", machine, reg.FindStringSubmatch(row)[0], temp)
						cluster.Node[machine] = strings.Split(reg.FindStringSubmatch(row)[0], "/")[0]
						found = true
					}
				}
			}
			if !found {
				log.Println("Machine up but no network...")
			} else {
				break
			}
		}
		time.Sleep(5 * time.Second)
		if inc >= 20 {
			log.Printf("May be check your target machine's availability.\nvirsh domifaddr %s --source agent still doesn't return a proper result...", machine)
		}
		inc++
	}
	return cluster, nil
}

func SeeIfOffline(cluster utils.Command, machine, remote string, silent bool) error {
	cmdstring := []string{"sudo", "virsh", "list", "--all"}
	var resp, err string
	if remote == "" {
		command := utils.SliceExec(cmdstring)
		if !silent {
			resp, err = utils.NiceBuffRunner(command, "/home/user")
		} else {
			out, err2 := command.CombinedOutput()
			resp = string(out)
			err = fmt.Sprintf("%s", err2)
		}
	} else {
		cmd := utils.SSHCommand(remote, cmdstring...)
		if !silent {
			resp, err = utils.NiceBuffRunner(cmd, "./")
		} else {
			out, err2 := cmd.CombinedOutput()
			resp = string(out)
			err = fmt.Sprintf("%s", err2)
		}
	}
	if err != "" {
		return errors.New(err)
	}
	for _, row := range strings.Split(resp, "\n") {
		if strings.Contains(row, machine) || strings.Contains(row, "shut off") {
			log.Printf("Machine % offline. Starting...\n", machine)
			StartVM(machine, remote)
		} else {
			log.Println("Domain unexistent. Removing all disks related to the machine...")
			delete(cluster.Node, machine)
			dir, _ := os.Getwd()
			cluster.SaveJSN(dir)
		}
	}
	return nil
}

func Destroy(cluster utils.Command, machine string) error {
	fmt.Printf("Destroying machine %v now...\n", machine)
	cmdstring1 := []string{"sudo", "virsh", "destroy", machine}
	cmdstring2 := []string{"sudo", "virsh", "undefine", machine}
	cmdstring3 := []string{"sudo", "virsh", "vol-delete", machine, "--pool", cluster.Pool.Name}
	command := utils.SliceExec(cmdstring1)
	_, err := utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	command = utils.SliceExec(cmdstring2)
	_, err = utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	command = utils.SliceExec(cmdstring3)
	_, err = utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	delete(cluster.Node, machine)
	return nil
}

func SanityCheck(cluster utils.Command) error {
	log.Println("Checking libvirt-based infra sanity...")
	cmdstring := []string{"sudo", "virsh", "pool-list"}
	command := utils.SliceExec(cmdstring)
	out, err := utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	for _, row := range strings.Split(out, "\n") {
		if strings.Contains(row, cluster.Pool.Name) {
			if !strings.Contains(row, "active") {
				cmdstring := []string{"sudo", "virsh", "pool-start", cluster.Pool.Name}
				command := utils.SliceExec(cmdstring)
				_, err := utils.NiceBuffRunner(command, "/home/user")
				if err != "" {
					return errors.New(err)
				}
			}
		}
	}
	cmdstring = []string{"sudo", "virsh", "net-list", "--all"}
	command = utils.SliceExec(cmdstring)
	out, err = utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	if cluster.Network.Name == "" {
		cluster.Network.Name = "default"
	}
	for _, row := range strings.Split(out, "\n") {
		if strings.Contains(row, cluster.Network.Name) {
			if strings.Contains(row, "inactive") {
				cmdstring := []string{"sudo", "virsh", "net-start", cluster.Network.Name}
				command := utils.SliceExec(cmdstring)
				_, err := utils.NiceBuffRunner(command, "/home/user")
				if err != "" {
					return errors.New(err)
				}
				cmdstring = []string{"sudo", "virsh", "net-autostart", cluster.Network.Name}
				command = utils.SliceExec(cmdstring)
				_, err = utils.NiceBuffRunner(command, "/home/user")
				if err != "" {
					return errors.New(err)
				}
			}
		}
	}
	return nil
}

func RefreshCluster(cluster utils.Command, flag bool) error {
	var err error
	for key, _ := range cluster.Node {
		cluster, err = CheckIfExists(cluster, key, "", true)
	}
	fmt.Printf("Machine                      IP\n-------------------------------------\n")
	for key, value := range cluster.Node {
		fmt.Printf("\n%s      %s\n", key, value)
	}
	return err
}
