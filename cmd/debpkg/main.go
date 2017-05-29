package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/xor-gate/debpkg"
)

var gOutputFile string
var gConfigFile string

func init() {
	flag.StringVar(&gConfigFile, "c", "", "Yaml configuration file")
	flag.StringVar(&gOutputFile, "o", "", "Debian output file")
	flag.Parse()
}

func main() {
	deb := debpkg.New()

	if gConfigFile == "" {
		gConfigFile = "debpkg.yml"
	}
	if err := deb.Config(gConfigFile); err != nil {
		log.Fatalf("Error while loading config file: %v", err)
	}
	if err := deb.Write(gOutputFile); err != nil {
		log.Fatalf("Error writing outputfile: %v", err)
		return
	}

	fmt.Println("debpkg: written:", gOutputFile)
}
