package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lextoumbourou/goodhosts"
)

const host = "api.asm.skype.com"

var regular bool

func main() {
	log.Printf("skyfka started")
	flag.BoolVar(&regular, "regular", false, "execute regularly")
	flag.Parse()
	if regular {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		// catch SIGINT, SYGTERM
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			s := <-sigs
			log.Printf("stopped with signal: %s", s)
			cancel()
		}()

		log.Printf("monitor %s", host)
		timer := time.NewTicker(time.Minute * 1)
	STOP:
		for {
			select {
			case <-timer.C:
				patchHosts()
			case <-ctx.Done():
				break STOP
			}
		}
	} else {
		log.Printf("patch host: %s", host)
		patchHosts()
	}
	log.Printf("skyfka stopped")
}

func patchHosts() {
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
