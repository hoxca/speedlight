package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	Log "github.com/apatters/go-conlog"
)

var i = true
var rootp string
var traversal filepath.WalkFunc = func(fp string, _ os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	path := filepath.ToSlash(fp)
	if i {
		rootp = path
	}
	i = false

	image := strings.TrimPrefix(path, rootp)
	/*
	   if runtime.GOOS == "windows" {
	           image = fmt.Sprintf("/%s", image)
	   }
	*/
	image = fmt.Sprintf("/%s", image)

	/*
	   Log.Debugf("path: %s\n", path)
	   Log.Debugf("rootp: %s\n", rootp)
	*/
	Log.Debugf("image: %s\n", image)

	regex := `/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*).*\.FIT`
	re := *regexp.MustCompilePOSIX(regex)
	splitline := re.FindAllStringSubmatch(image, -1)

	if len(splitline) == 1 {
		object := splitline[0][1]
		filter := splitline[0][2]
		expo, _ := strconv.Atoi(splitline[0][3])

		Log.Debugf("object %s", object)

		if !targetList.exist(object) {
			o = newTarget(object)
		}
		o.iterateFilter(filter, expo)
		targetList.set(object, o)
		if *verbosity == "debug" {
			o.printObject()
		}
	}
	return nil
}
