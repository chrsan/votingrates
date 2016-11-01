package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var verbose = flag.Bool("v", true, "verbose error output")

func main() {
	flag.Parse()
	var api HTTP
	regs, err := Regions(&api)
	if err != nil {
		die(err)
	}
	rates, err := Rates(&api, regs)
	if err != nil {
		die(err)
	}
	for _, r := range rates {
		fmt.Printf("%s %s %g%%\n", r.Year, strings.Join(r.Regs, ", "), r.Pct)
	}
}

func die(err error) {
	if *verbose {
		log.Fatalf("An unexpected error occurred: %s", err)
	} else {
		os.Exit(1)
	}
}
