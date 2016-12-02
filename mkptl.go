package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const resolver_count = 20
const channel_depth = 50

type Entry struct {
	name  string
	aux   string
	svc   int
	addrs []net.IP
	sub   *[]Entry
}

func (e *Entry) resolve() {
	addrs, aerr := net.LookupIP(e.name)

	if aerr == nil {
		e.addrs = addrs
	} else {
		log.Printf("A/AAAA resolution failed for %s: %s", e.name, string(aerr))
	}

	if e.svc == 80 {
		e.sub = make([]Entry)

		mxs, mxerr := net.LookupMX(e.name)
		nss, nserr := net.LookupNS(e.name)

		if mxerr == nil {
			e.mxnames = make([]string, len(mxs))
			for _, mx := range mxs {
				e.sub = append(e.sub, Entry{name: mx.Host, svc: 25})
			}
		} else {
			log.Printf("MX resolution failed for %s: %s", e.name, string(mxerr))
		}

		if nserr == nil {
			e.nsnames = make([]string, len(nss))
			for i, ns := range nss {
				e.sub = append(e.sub, Entry{name: ns.Host, svc: 53})
			}
		} else {
			log.Printf("NS resolution failed for %s: %s", e.name, string(nserr))
		}
	}
}

func main() {

	pending_entries = make(chan *Entry, channel_depth)
	finished_entries = make(chan *Entry, channel_depth)

	// start writing output
	go func() {
		for {
			e <- finished_entries
			if e == nil {
				break
			}
			for _, ip := range e.addrs {
				fmt.Fprintf(os.Stdout, "%s,%d,%s,%s",
					string(ip), e.svc, e.name, e.aux)
			}
		}
	}()

	// start resolving pending entries
	for i := 0; i < resolver_count; i++ {
		go func() {
			for {
				e <- pending_entries
				if e == nil {
					break
				}
				e.resolve()
				for _, se := range *e.sub {
					pending_entries <- &se
				}
				e.sub = nil
				finished_entries <- e
			}
		}()
	}

	// now scan input and convert it to unresolved entries
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line = scanner.Text()
		// comma: domain-name, aux pair
		// no comma: just a domain name
	}

}
