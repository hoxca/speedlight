package main

import (
	"flag"
	"log"
	"path/filepath"
)

var logFatal = log.Fatal

var (
	lightsdir   = flag.String("dir", "D:/Data/Voyager/Lights/", "lights directory")
	verbosity   = flag.String("level", "warn", "set log level of speedlight default warn")
	rotUnused   = flag.Bool("rotation", true, "tell if rotation is used")
	writeReport = flag.Bool("report", true, "write report to the filesystem")
)

func main() {
	writeConfig := writeDestination{true, *writeReport}

	flag.Parse()
	setUpLogs()

	err := filepath.Walk(*lightsdir, traversal)
	if err != nil {
		log.Println(err)
	}
	objectList.printObjects(writeConfig)
}
