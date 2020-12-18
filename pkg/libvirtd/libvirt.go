package libvirtd

import (
	"clusterer/pkg/data"
	"clusterer/pkg/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func CloneVol(cluster data.Command, seed string, machine string) error {
	log.Println("Clonning Volume(s)...")
	cmdstring := []string{"sudo", "virsh", "vol-clone", seed, machine, "--pool", cluster.Pool.Name}
	command := utils.SliceExec(cmdstring)
	_, err := utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	return nil
}

func CloneVM(cluster data.Command, seed string, machine string) error {
	log.Println("Clonning VM(s)...")
	cmdstring := []string{"sudo", "virt-clone", "-o", seed, "-n", machine, "--preserve-data", "-f", filepath.Join(cluster.Pool.Path, machine)}
	command := utils.SliceExec(cmdstring)
	_, err := utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	return nil
}

func StartVM(machine string) error {
	log.Println("Starting VM(s)...")
	cmdstring := []string{"sudo", "virsh", "start", machine}
	command := utils.SliceExec(cmdstring)
	_, err := utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	return nil
}

func WaitForIP(cluster data.Command, machine string) error {
	CheckIfExists(cluster, machine)
	return nil
}

func CheckIfExists(cluster data.Command, machine string) error {
	for {
		cmdstring := []string{"sudo", "virsh", "domifaddr", machine, "--source", "agent"}
		resp, err := exec.Command(cmdstring[0], cmdstring[1:]...).CombinedOutput()
		if err != nil {
			cmdstring2 := []string{"sudo", "virsh", "list"}
			resp2, err2 := exec.Command(cmdstring2[0], cmdstring2[1:]...).CombinedOutput()
			if err2 != nil {
				return err2
			}
			if strings.Contains(fmt.Sprintf("%s", string(resp2)), machine) {
				log.Println("Machine not ready yet...")
				time.Sleep(3 * time.Second)
			} else {
				log.Println("Looks like machine doesn't exist or offline...Will see if we remove it from cluster catalogue...")
				SeeIfOffline(cluster, machine)
			}
		} else {
			found := false
			for _, row := range strings.Split(fmt.Sprintf("%s", string(resp)), "\n") {
				if strings.Contains(row, "eth0") && strings.Contains(row, "ipv4") {
					//fmt.Printf("iaca IP: %s\n", strings.Split(row, " ")[len(strings.Split(row, " "))-1])
					cluster.Node[machine] = strings.Split(strings.Split(row, " ")[len(strings.Split(row, " "))-1], "/")[0]
					found = true
				}
			}
			if !found {
				log.Println("Machine up but no network...")
			} else {
				break
			}
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}

func SeeIfOffline(cluster data.Command, machine string) error {
	cmdstring := []string{"sudo", "virsh", "list", "--all"}
	command := utils.SliceExec(cmdstring)
	resp, err := utils.NiceBuffRunner(command, "/home/user")
	if err != "" {
		return errors.New(err)
	}
	for _, row := range strings.Split(resp, "\n") {
		if strings.Contains(row, machine) || strings.Contains(row, "shut off") {
			log.Printf("Machine % offline. Starting...\n", machine)
			StartVM(machine)
		} else {
			log.Println("Domain unexistent. Removing all disks related to the machine...")
			delete(cluster.Node, machine)
			dir, _ := os.Getwd()
			utils.SaveJSN(dir, cluster)
		}
	}
	return nil
}

func Destroy(cluster data.Command, machine string) error {
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

func SanityCheck(cluster data.Command) error {
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
