/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	//"fmt"
	"clusterer/pkg/utils"
	"clusterer/pkg/libvirtd"
	"github.com/spf13/cobra"
	"log"
)

// fakexenbmCmd represents the fakexenbm command
var (
	fakexenbmCmd = &cobra.Command{
	Use:   "fakexenbm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run3()
	},
}
remotehostIP = ""
//distro = ""
)

func init() {
	rootCmd.AddCommand(fakexenbmCmd)
	rootCmd.PersistentFlags().StringVar(&remotehostIP, "rmtip", "", "the remote host's IP address")
	//rootCmd.PersistentFlags().StringVar(&distro, "distro", "", "the SLES distribution")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fakexenbmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fakexenbmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


func run3() {
	
	path, err := utils.CopyRawDisks(remotehostIP, "10.84.149.229", distro) //CopyRawDisks(RemoteHostIP, FTPserverIP, Distro)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	machine, _ := libvirtd.DefineVMfromXML(path, "10.84.149.229")
	
	//machine := "sle15.2_fake_baremetal_xenvirthost_client"
	libvirtd.StartVM(machine, remotehostIP)
	//utils.ChangeXMLSpec(path)
}