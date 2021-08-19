package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func (cluster *Command) SumaServLogs(srvip, flag string) (err error) {
	switch flag {
	case "reposync":
		var out []byte
		cmd := []string{"ls", "-l", "--time-style=+\"%d %b %y %H:%M %z\"", "/var/log/rhn/reposync"}
		out, err = SSHCommand(srvip, cmd...).CombinedOutput()
		if err != nil {
			return
		}
		//fmt.Println(string(out))
		var period time.Duration
		var fileInfoRow string
		//fmt.Println(period)
		for _, row := range strings.Split(string(out), "\n") {
			if strings.Contains(row, ".log") {
				parsedTime := regexp.MustCompile(`\d{2}\ \w{3}\ \d{2}\ \d{2}:\d{2} .\d{4}`).FindString(string(row))
				//fmt.Println(parsedTime)
				day, err := time.Parse(time.RFC822Z, parsedTime)
				if err != nil {
					fmt.Println(err)
				}
				t := time.Now()
				if fmt.Sprintf("%v", period) == "0s" {
					period = t.Sub(day)
					fileInfoRow = row
				} else {
					if period > t.Sub(day) {
						period = t.Sub(day)
						fileInfoRow = row
					}
				}
				//fmt.Println(t.Sub(day))
			}
		}

		fileToTail := filepath.Join("/var/log/rhn/reposync", strings.Split(fileInfoRow, " ")[len(strings.Split(fileInfoRow, " "))-1])
		log.Printf("Looking into the file:   %s\n", fileToTail)
		//fmt.Println(fileToTail)
		//time.Sleep(100 * time.Second)
		cmd = []string{"tail", "-f", fileToTail}
		NiceBuffRunner(SSHCommand(srvip, cmd...), ".")
	}

	return
}
