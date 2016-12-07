package main

import (
	"bufio"
	"fmt"
	"github.com/mami-project/targets"
	"net"
	"os"
	"strings"
)

var priva, privb, privc *net.IPNet

func Is1918(ip net.IP) bool {
	return priva.Contains(ip) || privb.Contains(ip) || privc.Contains(ip)
}

func init() {
	_, priva, _ = net.ParseCIDR("10.0.0.0/8")
	_, privb, _ = net.ParseCIDR("172.16.0.0/12")
	_, privc, _ = net.ParseCIDR("192.168.0.0/16")
}

func main() {

	addrset := targets.MakeNameSet()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")

		ipstr := fields[0]
		port := fields[1]

		ip := net.ParseIP(ipstr)

		// skip non global unicast addresses
		if !ip.IsGlobalUnicast() {
			continue
		}

		// skip RFC 1918 addresses
		if Is1918(ip) {
			continue
		}

		if addrset.AddOnce(ipstr + "," + port) {
			fmt.Println(line)
		}

	}

}
