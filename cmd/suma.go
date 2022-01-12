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
	"clusterer/pkg/utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sclevine/agouti"
	"github.com/spf13/cobra"
)

// createchanCmd represents the createchan command
var (
	sumaCmd = &cobra.Command{
		Use:   "suma",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			createChanFromMU()
		},
	}
	url             string
	chanName        string
	srvIP           string
	orgCreds        bool
	createOrg       bool
	createChanRepo  bool
	deleteRepos     bool
	createReposJson string
	checkLogs       string
)

func init() {
	rootCmd.AddCommand(sumaCmd)
	rootCmd.PersistentFlags().StringVar(&srvIP, "srvip", "10.84.149.229", "SUMA Server's IP address")
	rootCmd.PersistentFlags().StringVar(&chanName, "channel", "default001", "custom channel's name")
	rootCmd.PersistentFlags().BoolVarP(&createOrg, "createorg", "o", false, "triggers org creation")
	rootCmd.PersistentFlags().BoolVar(&orgCreds, "addorgcred", false, "triggers org credentials creation")
	rootCmd.PersistentFlags().BoolVar(&deleteRepos, "deleterepos", false, "triggers all MU repos deletion")
	rootCmd.PersistentFlags().StringVar(&createReposJson, "reposjson", "", "path to json with all BV MU repositories")
	rootCmd.PersistentFlags().StringVar(&checkLogs, "checklogs", "", "checks the logs of ")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createchanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createchanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createChanFromMU() {

	if checkLogs != "" {
		Cluster.SumaServLogs(srvIP, checkLogs)
		os.Exit(3)
	}
	//time.Sleep(100 * time.Second)
	url := fmt.Sprintf("https://%s/", srvIP)
	page, err := utils.StartChromeSession(url, 300)
	utils.ErrHandler(err, "Error while starting chromium session:")

	if createOrg {
		log.Println("Creating organization...")
		CreateOrg(page)
		os.Exit(3)
	}

	if createReposJson != "" {
		log.Println("Creating Custom Channels from json...")
		Login(page)
		time.Sleep(3 * time.Second)
		a, err := utils.GetJson(createReposJson)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		channels, err := utils.CreateChannelNames(a)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		var counter int
		for index, value := range channels {
			urlChan := url + "rhn/channels/manage/Manage.do"
			err = page.Navigate(urlChan)
			err = CreateRepo(page, value, index, counter)
			if err != nil {
				log.Fatalf("Error: %v\n", err)
			}
			counter++
			if counter == len(channels) {
				os.Exit(3)
			}
		}

	}

	if deleteRepos {
		log.Println("Deleting Repositories...")
		Login(page)
		time.Sleep(3 * time.Second)
		for {
			urlChan := url + "rhn/channels/manage/Manage.do"
			err = page.Navigate(urlChan)
			err = DeleteAllRepos(page)
			if err != nil {
				log.Fatalf("Error: %v\n", err)
			}
		}
	}

	if orgCreds {
		log.Println("Adding Org Credentials...")
		Login(page)
		time.Sleep(4 * time.Second)
		OrgCredentials(page)
		time.Sleep(10 * time.Second)
	}

	if chanName != "default001" {
		log.Println("Creating Single Custom Channel...")
		//fmt.Println(fmt.Sprintf("custom_%v-%v-%v-%v-%v", time.Now().Month(), time.Now().Local().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second()))
		Login(page)
		time.Sleep(3 * time.Second)
		urlChan := url + "rhn/channels/manage/Manage.do"
		err = page.Navigate(urlChan)
		utils.ErrHandler(err, "Couldn't find page:")
		time.Sleep(3 * time.Second)
		CreateCustomChannel(page, urlChan, fmt.Sprintf("custom_%v-%v-%v-%v", time.Now().Local().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second()), chanName)
		time.Sleep(2 * time.Second)
	}
}

func FillElem(page *agouti.Page, patternToFind, fillingText string) {
	elem, err := utils.FindXpathElemOnPage(page, patternToFind, "simple")
	utils.ErrHandler(err, "Element wasn't found:")
	err = utils.FillForm(elem, fillingText)
	utils.ErrHandler(err, "Couldn't Fill the Form:")
}

func SelectFromForm(page *agouti.Page, patternToFind, selectOpt string) {
	elem := page.FindByName("prefix")
	//fmt.Println(elem)
	//utils.ErrHandler(err, "Element wasn't found:")
	err := elem.Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	//time.Sleep(30 * time.Second)
	err = elem.FindByXPath("option[@value=\"Dr.\"]").Click()
	utils.ErrHandler(err, fmt.Sprintf("Couldn't Find the Option %s:", selectOpt))
}

func ClickButton(page *agouti.Page, button string) {
	elem := page.FindByClass(button)
	//utils.ErrHandler(err, "Element wasn't found:")
	err := elem.Click()
	utils.ErrHandler(err, "Element couldn't be clicked: ")
}

func ClickLink(page *agouti.Page, link string) {
	elem, err := utils.FindXpathElemOnPage(page, link, "simple")
	utils.ErrHandler(err, "Element couldn't be found: ")
	err = elem.Click()
	utils.ErrHandler(err, "Element couldn't be clicked: ")
}

func OrgCredentials(page *agouti.Page) {
	err := page.FindByXPath("//ul[@class=\"level1\"]/li/div[@class=\" nodeLink \"]/a[@href=\"/rhn/admin/setup/ProxySettings.do\"]").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(1 * time.Second)
	//err = page.FindByXPath("//div[@class=\" leafLink \"]/a[text()=\"Organization Credentials\"]").Click()
	err = page.FirstByXPath("//a[@href=\"/rhn/admin/setup/MirrorCredentials.do\"]").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
}

func CreateOrg(page *agouti.Page) {
	FillElem(page, "#orgName", "SUSE")
	FillElem(page, "#loginname", "admin")
	FillElem(page, "#desiredpass", "admin")
	FillElem(page, "#confirmpass", "admin")
	FillElem(page, "#email", "admin@admin.com")
	SelectFromForm(page, "#prefix", "Miss")
	FillElem(page, "#firstNames", "Admin")
	FillElem(page, "#lastName", "Adminsson")
	ClickButton(page, "btn-success")
}

func Login(page *agouti.Page) *agouti.Page {
	FillElem(page, "#username-field", "admin")
	FillElem(page, "#password-field", "admin")
	ClickButton(page, "btn-success")
	return page
}

func FindParentChannel(page *agouti.Page, url, channelName, MUrepo string) {

}

func CreateRepo(page *agouti.Page, label, url string, count int) (err error) {
	err = page.FindByXPath("//a[@href=\"/rhn/channels/manage/repos/RepoList.do\"]").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(3 * time.Second)
	err = page.FindByXPath("//a[text()=\"Create Repository\"]").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(3 * time.Second)
	err = page.FindByName("label").Fill(label)
	err = page.FindByName("url").Fill(url)
	ClickButton(page, "btn-success")
	time.Sleep(3 * time.Second)
	return
}

func DeleteAllRepos(page *agouti.Page) (err error) {
	err = page.FindByXPath("//a[@href=\"/rhn/channels/manage/repos/RepoList.do\"]").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(1 * time.Second)
	err = page.FirstByXPath("//tr[@class=\"list-row-odd\"]/td/a").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(2 * time.Second)
	err = page.FindByXPath("//i[@title=\"Delete Repository\"]").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(2 * time.Second)
	err = page.FindByXPath("//input[@value=\"Delete Repository\"]").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	return
}

func CreateCustomChannel(page *agouti.Page, url, channelName, MUrepo string) {
	/*
		err := page.FindByClass("fa-plus").Click()
		utils.ErrHandler(err, "Element can't be clicked:")
		FillElem(page, "#name", "admin001")
		FillElem(page, "#label", "kuka001")
		FillElem(page, "#summary", "bujumbura")
		ClickButton(page, "btn-success")
		time.Sleep(2 * time.Second)
	*/
	//elem := page.Find("admin001")
	////*[@id="spacewalk-content"]/div[2]/ul/li[5]/a
	//template := fmt.Sprintf("//tr[@class=\"list-row-odd\"]/td[1]/a", channelName)

	//elem := page.FirstByXPath("//tr[@class=\"list-row-odd\"]/td[1]/a").Click()
	//span[text()='All Sector ETFs']
	//elem := page.FirstByXPath("//tr[@class=\"list-row-odd\"]")
	formatForChannel := fmt.Sprintf("//a[text()=\"%s\"]", channelName)
	err := page.FindByXPath(formatForChannel).Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(2 * time.Second)
	elem := page.FindByClass("spacewalk-content-nav")
	err = elem.FindByXPath("//ul/li[5]/a").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(2 * time.Second)
	elem = page.FindByClass("action-button-wrapper")
	err = elem.FindByClass("btn-default").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(2 * time.Second)
	err = page.FindByName("label").Fill("mu-20208")
	err = page.FindByName("url").Fill(MUrepo)
	ClickButton(page, "btn-success")
	time.Sleep(2 * time.Second)
	//fmt.Println(page.URL())
	elem = page.FindByClass("spacewalk-content-nav")
	err = elem.FindByXPath("//ul[@class=\"nav nav-tabs nav-tabs-pf\"]/li[2]/a").Click()
	///html/body/div[2]/div/div/section/div[3]/ul[2]/li[2]/a
	//err = elem.FindByXPath("//li[2]/a").Click()
	utils.ErrHandler(err, "Element can't be clicked:")
	time.Sleep(1 * time.Second)
	ClickButton(page, "btn-success")
}
