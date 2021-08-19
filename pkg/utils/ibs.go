package utils

import (
	"fmt"
	"regexp"
)

func SyncMUChannel(mu string) error {
	var maintUpdate MU
	mu = "SUSE:Maintenance:20223:244004"
	pref := regexp.MustCompile(`\W{4}:W{11}`)
	maintUpdate.Prefix = pref.FindString(mu)
	fmt.Println(pref)
	return nil
}
