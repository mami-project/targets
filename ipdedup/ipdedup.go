package main

import (
	"bufio"
	"fmt"
	"github.com/mami-project/targets"
	"os"
	"strings"
)

func main() {

	addrset := targets.MakeNameSet()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")

		ip := fields[0]
		port := fields[1]

		if addrset.AddOnce(ip + "," + port) {
			fmt.Println(line)
		}

	}

}
