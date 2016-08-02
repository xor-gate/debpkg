package main

import (
	"fmt"
	"flag"
	"github.com/xor-gate/debpkg"
)

var g_outputFile string
var g_configFile string

func init() {
	flag.StringVar(&g_configFile, "c", "", "Yaml configuration file")
	flag.StringVar(&g_outputFile, "o", "", "Debian output file")
	flag.Parse()
}

func main() {
	deb := debpkg.New()

	if g_configFile != "" {
		err := deb.Config(g_configFile)
		if err != nil {
			fmt.Println("Error while loading config file", g_configFile)
			return
		}
	}

	deb.Write(g_outputFile)
	fmt.Println("Written", g_outputFile)
}
