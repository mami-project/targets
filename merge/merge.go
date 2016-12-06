package main

import (
	"bufio"
	"fmt"
	"github.com/mami-project/targets"
	"log"
	"os"
	"strings"
)

func main() {
	var files []*os.File
	var scanners []*bufio.Scanner

	nameset := targets.MakeNameSet()

	for i := 1; i < len(os.Args); i++ {
		f, err := os.Open(os.Args[i])
		if err != nil {
			log.Fatal("couldn't open %s: %s", os.Args[i], err.Error())
		}
		files = append(files, f)

		s := bufio.NewScanner(f)
		s.Split(bufio.ScanLines)
		scanners = append(scanners, s)
	}

	for {
		stillScanning := 0
		for _, s := range scanners {
			if s.Scan() {
				if nameset.AddOnce(strings.ToLower(s.Text())) {
					fmt.Println(s.Text())
				}
				stillScanning++
			}
		}
		if stillScanning == 0 {
			break
		}
	}
}
