package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Entry struct {
	name  string
	aux   string
	svc   int
	addrs []net.IP
	sub   *[]Entry
	do_mx bool
	do_ns bool
}

func (e *Entry) resolve() {
	addrs, aerr := net.LookupIP(e.name)

	if aerr == nil {
		e.addrs = addrs
	} else {
		e.addrs = make([]net.IP, 1)
	}

	if e.do_ns || e.do_mx {
		e.sub = new([]Entry)

		if e.do_ns {
			nss, nserr := net.LookupNS(e.name)
			if nserr == nil {
				for _, ns := range nss {
					*e.sub = append(*e.sub, Entry{name: ns.Host, svc: 53})
				}
			}
		}

		if e.do_mx {
			mxs, mxerr := net.LookupMX(e.name)

			if mxerr == nil {
				for _, mx := range mxs {
					*e.sub = append(*e.sub, Entry{name: mx.Host, svc: 25})
				}
			}
		}
	}
}

func do_resolution(e *Entry, finished chan *Entry, limiter chan struct{}, wait *sync.WaitGroup) {
	wait.Add(1)
	limiter <- struct{}{}
	e.resolve()
	if e.sub != nil {
		for _, se := range *e.sub {
			go do_resolution(&se, finished, limiter, wait)
		}
	}
	finished <- e
	_ = <-limiter
	wait.Done()
}

func main() {

	// command-line flags
	default_svc := flag.Int("svc", 80, "Port number for top-level resolutions")
	do_mx := flag.Bool("mx", false, "Also attempt to resolve MX addresses")
	do_ns := flag.Bool("ns", false, "Also attempt to resolve NS addresses")
	also := flag.String("also", "", "Also attempt to resolve additional name within domain")
	resolver_count := flag.Int("resolvers", 32, "Maximum concurrent resolutions")
	flag.Parse()

	// some channels
	finished := make(chan *Entry, 32)
	limiter := make(chan struct{}, *resolver_count)
	resolver_wait := new(sync.WaitGroup)
	output_wait := new(sync.WaitGroup)

	// start writing output
	go func() {
		var e *Entry
		output_wait.Add(1)
		for {
			e = <-finished
			if e == nil {
				break
			}
			for _, ip := range e.addrs {
				fmt.Fprintf(os.Stdout, "%s,%d,%s,%s\n",
					ip.String(), e.svc, e.name, e.aux)
			}
		}
		output_wait.Done()
	}()

	// now scan input and convert it to unresolved entries
	lineNum := 0
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		fields := strings.Split(line, ",")

		if len(fields) >= 1 {
			e := new(Entry)
			e.name = fields[0]
			e.svc = *default_svc
			e.do_mx = *do_mx
			e.do_ns = *do_ns
			if len(fields) >= 2 {
				e.aux = fields[1]
			} else {
				e.aux = strconv.Itoa(lineNum)
			}

			go do_resolution(e, finished, limiter, resolver_wait)

			if len(*also) > 0 {
				we := new(Entry)
				we.name = *also + "." + fields[0]
				we.svc = *default_svc
				we.aux = e.aux

				go do_resolution(we, finished, limiter, resolver_wait)
			}
		}
	}

	resolver_wait.Wait()
	close(finished)
	output_wait.Wait()
}
