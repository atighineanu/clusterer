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
	"time"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "starts newly created machines in the cluster",
		Long: `starts the created machines in the cluster;
		it waits for the VM IPs to be up, before reporting that all machines are ready`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting cluster nodes...")
			StartCluster()
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	//fmt.Printf("%+v", cluster)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartCluster() {
	cluster, err := utils.OpenJSN(filepath.Join(RootDir, "cluster.json"))
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
	libvirtd.SanityCheck(*cluster)
	for index, _ := range cluster.Node {
		libvirtd.StartVM(index)
	}
	time.Sleep(17 * time.Second)
	for index, _ := range cluster.Node {
		libvirtd.WaitForIP(*cluster, index)
		time.Sleep(3 * time.Second)
	}
	if err := utils.SaveJSN(RootDir, *cluster); err != nil {
		log.Printf("Error while saving the json file: %v", err)
	}
}
