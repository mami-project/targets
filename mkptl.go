package main

import (
	"log"
	"net"
)

type Entry struct {
	name    string
	addrs   []net.IP
	mxnames []string
	nsnames []string
}

func (e *Entry) resolve() {

	addrs, aerr := net.LookupIP(e.name)
	mxs, mxerr := net.LookupMX(e.name)
	nss, nserr := net.LookupNS(e.name)

	if aerr == nil {
		e.addrs = addrs
	}

	if mxerr == nil {
		e.mxnames = make([]string, len(mxs))
		for i, mx := range mxs {
			e.mxnames[i] = mx.Host
		}
	}

	if nserr == nil {
		e.mxnames = make([]string, len(mxs))
		for i, mx := range mxs {
			e.mxnames[i] = mx.Host
		}
	}
}

func main() {

}
