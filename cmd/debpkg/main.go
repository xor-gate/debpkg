package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xor-gate/debpkg"
)

var (
	outputFile    string
	configFile    string
	versionNumber string
)

func init() {
	flag.StringVar(&configFile, "c", "debpkg.yml",
		"YAML configuration file")
	flag.StringVar(&outputFile, "o", "",
		"Debian output file")
	flag.StringVar(&versionNumber, "v", os.Getenv("DEBPKG_VERSION"),
		"Package version number (or via DEBPKG_VERSION environment variable)")
	flag.Parse()
}

func main() {
	deb := debpkg.New()
	if err := deb.Config(configFile); err != nil {
		log.Fatalf("Error while loading config file: %v", err)
	}
	if versionNumber != "" {
		deb.SetVersion(versionNumber)
	}
	if err := deb.Write(outputFile); err != nil {
		log.Fatalf("Error writing outputfile: %v", err)
	}
	fmt.Println("debpkg: written:", outputFile)
}
