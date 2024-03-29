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
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var (
	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "destroys all the machines in the cluster.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			Run()
		},
	}
	all = false
)

func init() {
	rootCmd.AddCommand(destroyCmd)
	rootCmd.PersistentFlags().BoolVar(&all, "all", false, "delete all nodes and resources")
}

func Run() {
	Cluster, err := utils.OpenJSN(filepath.Join(RootDir, "cluster.json"))
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
	if all {
		for index, _ := range Cluster.Node {
			libvirtd.Destroy(*Cluster, index)
		}
	}
	if err := Cluster.SaveJSN(RootDir); err != nil {
		log.Printf("Error while saving the json file: %v", err)
	}
}
