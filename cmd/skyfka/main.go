package main

import (
	"log"
	"net"
	"strings"

	"github.com/lextoumbourou/goodhosts"
)

func main() {
	host := "api.asm.skype.com"
	addrs, err := net.LookupHost(host)
	if err != nil {
		log.Fatalf("failed to lookup: %s", err)
	}
	filtered := getRemoteAddrs(addrs)
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		log.Fatalf("failed to open hosts file: %s", err)
	}
	if ok := hosts.IsWritable(); !ok {
		log.Fatalf("hosts file is not writable")
	}
	for _, addr := range filtered {
		if err := hosts.Add(addr, host); err != nil {
			log.Printf("add record into hosts file: %s", err)
		}
	}
	if err := hosts.Flush(); err != nil {
		log.Fatal(err)
	}
}

func getRemoteAddrs(input []string) []string {
	filtered := []string{}
	for _, addr := range input {
		if !isLocal(addr) {
			filtered = append(filtered, addr)
		}
	}
	return filtered
}

func isLocal(addr string) bool {
	for _, local := range []string{"::", "127.0.0.1", "localhost"} {
		if strings.HasPrefix(addr, local) {
			return true
		}
	}
	return false
}
