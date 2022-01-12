/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"clusterer/pkg/libvirtd"
	"clusterer/pkg/utils"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

// spawnCmd represents the spawn command

var (
	spawnCmd = &cobra.Command{
		Use:   "spawn",
		Short: "command to create machines; works with --workers --masters --pool --stackname (and other) flags",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Printf("here is amount of workers: %v\n", workers)
			setup()
			Deploy()
		},
	}
	deploy   = "libvirt"
	stack    = ""
	pool     = "default"
	poolpath = "/home/qcows"
	distro   = ""
	workers  = 0
	masters  = 0
	jsn      = false
	err      error
)

func init() {
	rootCmd.AddCommand(spawnCmd)
	rootCmd.PersistentFlags().BoolVar(&jsn, "json", false, "triggers output in json")
	rootCmd.PersistentFlags().StringVar(&stack, "stackname", "default", "number of workers in the cluster")
	rootCmd.PersistentFlags().IntVar(&workers, "workers", 0, "number of workers in the cluster")
	rootCmd.PersistentFlags().IntVar(&masters, "masters", 0, "number of masters in the cluster")
	rootCmd.PersistentFlags().StringVarP(&pool, "pool", "p", "default", "name of pool for the project")
	rootCmd.PersistentFlags().StringVarP(&distro, "distro", "d", "", "name of distro in the cluster")
	//rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	//rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	//viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	Cluster, err = utils.OpenJSN(filepath.Join(RootDir, "cluster.json"))
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
}

func setup() {
	//fmt.Println(Cluster)
	Cluster.Deploy = deploy
	Cluster.StackName = stack
	Cluster.Workers.Count = workers
	Cluster.Masters.Count = masters
	if pool != "" {
		Cluster.Pool.Name = pool
		Cluster.Pool.Path = poolpath
	}
	Cluster.SeedVol_Leap = "opensuse-seed.qcow2"
	Cluster.SeedVM_Leap = "opensuse-seed"
	Cluster.SeedVM_Ubuntu = "ubuntu_21.04-seed"
	Cluster.SeedVol_Ubuntu = "ubuntu_21.04-seed.qcow2"
	if distro != "" {
		Cluster.Workers.Distro = distro
		Cluster.Masters.Distro = distro
	}
	Cluster.Node = make(map[string]string)
	if err := Cluster.SaveJSN(RootDir); err != nil {
		log.Printf("Error while saving the json file: %v", err)
	}
}

func Deploy() {
	for i := 0; i < Cluster.Workers.Count; i++ {
		Cluster.Node[fmt.Sprintf("%s-%s-%v", Cluster.StackName, "workers", i)] = ""
	}
	for i := 0; i < Cluster.Masters.Count; i++ {
		Cluster.Node[fmt.Sprintf("%s-%s-%v", Cluster.StackName, "masters", i)] = ""
	}
	for index, _ := range Cluster.Node {
		switch distro {
		case "leap":
			libvirtd.CloneVol(*Cluster, Cluster.SeedVol_Leap, index)
			libvirtd.CloneVM(*Cluster, Cluster.SeedVM_Leap, index)
		case "ubuntu":
			libvirtd.CloneVol(*Cluster, Cluster.SeedVol_Ubuntu, index)
			libvirtd.CloneVM(*Cluster, Cluster.SeedVM_Ubuntu, index)
		}

	}
	if err := Cluster.SaveJSN(RootDir); err != nil {
		log.Printf("Error while saving the json file: %v", err)
	}
}
