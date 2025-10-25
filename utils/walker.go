package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	Log "github.com/apatters/go-conlog"
)

var i = true
var rootp string
var now = time.Now()

var Flatsversal filepath.WalkFunc = func(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	path := filepath.ToSlash(fp)
	if i {
		rootp = path
	}
	i = false
	// layout := "2006-01-02 15:04:05 +0200 CEST"

	var cutoff = 48 * time.Hour

	image := strings.TrimPrefix(path, rootp)
	if runtime.GOOS == "windows" {
		image = fmt.Sprintf("/%s", image)
	} else {

		image = fmt.Sprintf("%s", image)
	}

	if !fi.IsDir() && fi.Name() != ".DS_Store" {
		if diff := now.Sub(fi.ModTime()); diff < cutoff {

			Log.Debugf("image: %s", image)

			regex := `/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*)s_BIN1_(.*)C_GA.*_[[:digit:]]{8}_[[:digit:]]{6}_[[:digit:]]{3}_PA([[:digit:]]{3}\.[[:digit:]]{2})_[EW]\.FIT`
			re := *regexp.MustCompilePOSIX(regex)
			splitline := re.FindAllStringSubmatch(image, -1)

			if len(splitline) == 1 {
				object := splitline[0][1]
				filter := splitline[0][2]
				expo, _ := strconv.Atoi(splitline[0][3])
				temperature, _ := strconv.Atoi(splitline[0][4])
				rotval, _ := strconv.ParseFloat(splitline[0][5], 32)
				rotation := float32(rotval)

				if !RotUsed {
					rotation = 666
				}

				Log.Debugf("object %s filter %s expo %d temperature %d rotation %.2f", object, filter, expo, temperature, rotation)
				addFlats(rotation, filter)

			}
			return nil
		}
	}

	return nil
}

var Traversal filepath.WalkFunc = func(fp string, _ os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	// fmt.Printf("%q\n", fi.ModTime())
	path := filepath.ToSlash(fp)
	if i {
		rootp = path
	}
	i = false

	image := strings.TrimPrefix(path, rootp)
	if runtime.GOOS == "windows" {
		image = fmt.Sprintf("/%s", image)
	} else {

		image = fmt.Sprintf("%s", image)
	}
	Log.Debugf("image: %s\n", image)

	// Initial Kosmodrom file pattern: CR399_LIGHT_L_600s_BIN1_-20C_002_20221015_222558_734_E .

	regex := `/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*)s_BIN1_(.*)C_[[:digit:]]{3}_[[:digit:]]{8}_[[:digit:]]{6}_([[:digit:]]{3})_[EW]*.*\.FIT`
	re := *regexp.MustCompilePOSIX(regex)
	splitline := re.FindAllStringSubmatch(image, -1)

	if len(splitline) == 1 {
		object := splitline[0][1]
		filter := splitline[0][2]
		expo, _ := strconv.Atoi(splitline[0][3])
		temperature, _ := strconv.Atoi(splitline[0][4])
		rotation, _ := strconv.Atoi(splitline[0][5])
		if !RotUsed {
			rotation = 666
		}

		Log.Debugf("object %s filter %s expo %d temperature %d", object, filter, expo, temperature)

		target := fmt.Sprintf("%s~%d~%d", object, temperature, expo)
		o = addTarget(target, rotation)
		o.addFilter(target, filter)
	}
	return nil
}
