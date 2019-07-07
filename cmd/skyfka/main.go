package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bogdanovich/dns_resolver"
	"github.com/lextoumbourou/goodhosts"
)

const defaultHost = "api.asm.skype.com"

var regular bool
var host string

func main() {
	log.Printf("skyfka started")
	flag.BoolVar(&regular, "regular", false, "execute regularly")
	flag.StringVar(&host, "host", defaultHost, "skype api host")
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
				patchHosts(host)
			case <-ctx.Done():
				break STOP
			}
		}
	} else {
		log.Printf("patch host: %s", host)
		patchHosts(host)
	}
	log.Printf("skyfka stopped")
}

func lookupHost(target string) ([]net.IP, error) {
	resolver := dns_resolver.New([]string{"8.8.8.8", "8.8.4.4"})
	resolver.RetryTimes = 5

	ip, err := resolver.LookupHost(target)
	if err != nil {
		return nil, err
	}
	return ip, nil
}

func patchHosts(host string) {
	addrs, err := lookupHost(host)
	// addrs, err := net.LookupHost(host)
	if err != nil {
		log.Fatalf("failed to lookup: %s", err)
	}
	log.Printf("IPs: %v", addrs)
	filtered := getRemoteAddrs(addrs)
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		log.Fatalf("failed to open hosts file: %s", err)
	}
	if ok := hosts.IsWritable(); !ok {
		log.Fatalf("hosts file is not writable")
	}

	for _, addr := range filtered {
		if err := hosts.Add(addr.String(), host); err != nil {
			log.Printf("add record into hosts file: %s", err)
		}
	}

	if err := hosts.Flush(); err != nil {
		log.Fatal(err)
	}
}

func getRemoteAddrs(input []net.IP) []net.IP {
	filtered := []net.IP{}
	for _, addr := range input {
		if !addr.IsLoopback() {
			filtered = append(filtered, addr)
		}
	}
	return filtered
}
