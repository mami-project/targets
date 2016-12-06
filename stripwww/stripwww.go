package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "www.") == 0 {
			line = strings.Replace(line, "www.", "", 1)
		}
		fmt.Println(line)
	}
}
