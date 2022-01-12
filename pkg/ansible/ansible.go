package ansible

import (
	"clusterer/pkg/data"
	"clusterer/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type InvTempStruct struct {
	StackName string
	AllIps    string
	IpMasters string
	IpWorkers string
}

func PopulateInventory(cluster utils.Command, rootdir string) error {
	var temp InvTempStruct
	temp.StackName = strings.ToLower(cluster.StackName)
	for key, value := range cluster.Node {
		if strings.Contains(strings.ToLower(key), "master") {
			temp.IpMasters += value + "\n"
		}
		if strings.Contains(strings.ToLower(key), "worker") {
			temp.IpWorkers += value + "\n"
		}
	}
	temp.AllIps = temp.IpMasters + temp.IpWorkers
	if err := ExecTemplate(temp, rootdir); err != nil {
		return err
	}
	return nil
}

func ExecTemplate(temp InvTempStruct, rootdir string) error {
	var err error
	var f *os.File

	f, err = os.Create(filepath.Join(rootdir, "ansible/inventory"))
	if err != nil {
		log.Fatalf("couldn't create the file...%s", err)
	}

	templ, err := template.New("test").Parse(data.Inventorytmpl)
	if err != nil {
		return err
	}
	err = templ.Execute(f, temp)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}
