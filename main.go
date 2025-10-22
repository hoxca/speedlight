package main

import (
	"flag"
	"log"
	"path/filepath"
)

// var logFatal = log.Fatal

var (
	lightsdir = flag.String("dir", "D:/Data/Voyager/Lights/", "lights directory")
	verbosity = flag.String("level", "warn", "set log level of speedlight default warn")
	rotUnused = flag.Bool("rot", true, "tell if rotation is used")
)

func main() {
	// writeConfig := writeDestination{true, true} .
	flag.Parse()
	setUpLogs()

	err := filepath.Walk(*lightsdir, traversal)
	if err != nil {
		log.Println(err)
	}
	objectList.printObjects()
	// targetList.printObjects(writeConfig) .
}
