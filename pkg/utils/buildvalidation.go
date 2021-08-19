package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetJson() (a map[string]map[string]string, err error) {
	a = make(map[string]map[string]string)
	file, err := os.Open("bv.json")
	if err != nil {
		return
	}
	bytechunk, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytechunk, &a)
	if err != nil {
		return
	}
	//fmt.Printf("%+v\n", a)
	return
}

func CreateChannelNames(a map[string]map[string]string) (chans map[string]string, err error) {
	channels := make(map[string]map[string][]string)
	chans = make(map[string]string)
	packDistro := make(map[string][]string)
	for index, value := range a {
		for ind, val := range value {
			if channels[val] != nil {
				packDistro = channels[val]
			}
			packDistro[ShortenRepoLabels(ind)] = append(packDistro[ShortenRepoLabels(ind)], ShortenRepoLabels(index))
			channels[val] = packDistro
			packDistro = make(map[string][]string)
		}
	}
	var chanShortName string
	for ind, val := range channels {
		for index, elem := range val {
			if index == "" {
				index = "0"
			}
			chanShortName = fmt.Sprintf("%s-%s", index, strings.Join(elem, "_"))
			chans[ind] = chanShortName
			//fmt.Printf("%s    %s\n", ind, chanShortName)
		}
	}
	//time.Sleep(100 * time.Second)
	return
}

func ShortenRepoLabels(distro string) (shortDistro string) {
	dstr := strings.Split(distro, "_")
	for i := 0; i < len(dstr); i++ {
		switch dstr[i] {
		case "prometheus":
			dstr[i] = "prm"
		case "basesystem":
			dstr[i] = "bs"
		case "traditional":
			dstr[i] = "tr"
		case "salt":
			dstr[i] = "sl"
		case "salt2":
			dstr[i] = "sl2"
		case "server":
			dstr[i] = "s"
		case "proxy":
			dstr[i] = "p"
		case "sle11sp4":
			dstr[i] = "11.4"
		case "sle12sp4":
			dstr[i] = "12.4"
		case "sle12sp5":
			dstr[i] = "12.5"
		case "sle15":
			dstr[i] = "15"
		case "sle15sp1":
			dstr[i] = "15.1"
		case "sle15sp2":
			dstr[i] = "15.2"
		case "sle15sp3":
			dstr[i] = "15.3"
		case "ubuntu1804":
			dstr[i] = "u1804"
		case "ubuntu2004":
			dstr[i] = "u2004"
		case "client":
			dstr[i] = "c"
		case "minion":
			dstr[i] = "m"
		}
	}
	if len(dstr) == 1 {
		shortDistro = dstr[0]
	} else {
		if len(dstr) == 2 {
			shortDistro = fmt.Sprintf("%s_%s", dstr[0], dstr[1])
		}
	}
	return
}
