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

type NameSet struct {
	names map[string]struct{}
	lock  sync.RWMutex
}

func (t *NameSet) addOnce(name string) bool {
	t.lock.RLock()
	_, ok := t.names[name]
	t.lock.RUnlock()

	if ok {
		return false
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	_, ok = t.names[name]
	if ok {
		return false
	}
	t.names[name] = struct{}{}
	return true
}

func NewNameSet() *NameSet {
	out := new(NameSet)
	out.names = make(map[string]struct{})
	return out
}

func (e *Entry) resolve() {
	addrs, aerr := net.LookupIP(e.name)

	if aerr == nil {
		e.addrs = addrs
	} else {
		e.addrs = make([]net.IP, 0)
	}

	if e.do_ns || e.do_mx {
		e.sub = new([]Entry)

		if e.do_ns {
			nss, nserr := net.LookupNS(e.name)
			if nserr == nil {
				for _, ns := range nss {
					*e.sub = append(*e.sub, Entry{name: ns.Host, svc: 53, aux: e.aux})
				}
			}
		}

		if e.do_mx {
			mxs, mxerr := net.LookupMX(e.name)

			if mxerr == nil {
				for _, mx := range mxs {
					*e.sub = append(*e.sub, Entry{name: mx.Host, svc: 25, aux: e.aux})
				}
			}
		}
	}
}

func do_resolution(
	e *Entry,
	finished chan *Entry,
	limiter chan struct{},
	duplicates *NameSet,
	wait *sync.WaitGroup) {

	wait.Add(1)
	defer wait.Done()

	if duplicates.addOnce(strings.ToLower(e.name)) {
		limiter <- struct{}{}
		e.resolve()
		if e.sub != nil {
			for _, se := range *e.sub {
				go do_resolution(&se, finished, limiter, duplicates, wait)
			}
		}
		finished <- e
		_ = <-limiter
	}
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

	// wait groups
	resolver_wait := new(sync.WaitGroup)
	output_wait := new(sync.WaitGroup)

	// duplicate table
	duplicates := NewNameSet()

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

			go do_resolution(e, finished, limiter, duplicates, resolver_wait)

			if len(*also) > 0 {
				we := new(Entry)
				we.name = *also + "." + fields[0]
				we.svc = *default_svc
				we.aux = e.aux

				go do_resolution(we, finished, limiter, duplicates, resolver_wait)
			}
		}
	}

	resolver_wait.Wait()
	close(finished)
	output_wait.Wait()
}
