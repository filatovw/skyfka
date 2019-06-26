package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	hostsFile string
)

func main() {
	flag.StringVar(&hostsFile, "hosts-file", "", "path to hosts file on your system")
	flag.Parse()

	addrs, err := net.LookupHost("api.asm.skype.com")
	if err != nil {
		log.Fatalf("failed to lookup: %s", err)
	}
	filtered := getRemoteAddrs(addrs)

	fmt.Printf("%v", filtered)
}

func getRemoteAddrs(input []string) []string {
	filtered := []string{}
	for _, addr := range input {
		if strings.HasPrefix(addr, "::") {
			continue
		}
		filtered = append(filtered, addr)
	}
	return filtered
}

func updateHosts(pathToHosts string, addrs []string) (bool, error) {
	f, err := os.Open(pathToHosts)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close %s", pathToHosts)
		}
	}()
	return true, nil
}
