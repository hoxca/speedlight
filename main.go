package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var i = true
var rootp string

type target struct {
	name  string
	L     int
	Lexpo int
	R     int
	Rexpo int
	G     int
	Gexpo int
	B     int
	Bexpo int
	S     int
	Sexpo int
	H     int
	Hexpo int
	O     int
	Oexpo int
}

type targets map[string]target

var targetList = targets{}
var o *target

type writeDestination struct {
	writeToConsole bool
	writeToFile    bool
}

var traversal filepath.WalkFunc = func(fp string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	path := filepath.ToSlash(fp)
	if i {
		rootp = path
	}
	i = false

	image := strings.TrimPrefix(path, fmt.Sprintf("%s", rootp))
	if runtime.GOOS == "windows" {
		image = fmt.Sprintf("/%s", image)
	}
	/*
		fmt.Printf("path: %s\n", path)
		fmt.Printf("rootp: %s\n", rootp)
		fmt.Printf("image: %s\n", image)
	*/

	regex := `/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*).*\.FIT`
	re := *regexp.MustCompilePOSIX(regex)
	splitline := re.FindAllStringSubmatch(image, -1)

	if len(splitline) == 1 {
		object := splitline[0][1]
		filter := splitline[0][2]
		expo, _ := strconv.Atoi(splitline[0][3])

		if !targetList.exist(object) {
			o = newTarget(object, expo)
		}
		o.iterateFilter(filter, expo)

		targetList.set(object, o)
	}
	return nil
}

var lightsdir = flag.String("dir", "D:/Data/Voyager/Lights/", "lights directory")

func main() {

	writeConfig := writeDestination{true, true}
	flag.Parse()

	err := filepath.Walk(*lightsdir, traversal)
	if err != nil {
		log.Println(err)
	}
	targetList.printObjects(writeConfig)
}

func (ts targets) get(key string) *target {
	obj := ts[key]
	return &obj
}

func (ts targets) set(key string, value *target) {
	ts[key] = *value
}

func (ts targets) exist(key string) bool {
	for k := range ts {
		if k == key {
			return true
		}
	}
	return false
}

func (t *target) printObject() {
	fmt.Printf("filters: %d %d %d %d %d %d %d\n", t.L, t.R, t.G, t.B, t.S, t.H, t.O)
	fmt.Printf("expt. %d %d %d %d %d %d %d\n", t.Lexpo, t.Rexpo, t.Gexpo, t.Bexpo, t.Sexpo, t.Hexpo, t.Oexpo)
}

func (t *target) iterateFilter(filter string, expo int) {
	switch filter {
	case "L":
		t.L++
		t.Lexpo = expo
	case "R":
		t.R++
		t.Rexpo = expo
	case "G":
		t.G++
		t.Gexpo = expo
	case "B":
		t.B++
		t.Bexpo = expo
	case "S":
		t.S++
		t.Sexpo = expo
	case "H":
		t.H++
		t.Hexpo = expo
	case "O":
		t.O++
		t.Oexpo = expo
	}
}

func newTarget(object string, expo int) *target {
	return &target{
		name: object,
		L:    0,
		R:    0,
		G:    0,
		B:    0,
		S:    0,
		H:    0,
		O:    0,
	}
	//	return &newobject
}

func (ts *targets) getTargets() []string {
	ret := []string{}
	for _, v := range *ts {
		ret = append(ret, v.name)
	}
	return ret
}

func (ts *targets) printObjects(wdest writeDestination) {

	targets := ts.getTargets()

	if wdest.writeToConsole {
		fmt.Printf("Targets list: %q\n\n", targets)
		for _, v := range *ts {
			fmt.Printf("Object: %-30s %s%s\n", v.name, "Total: ", secondsToHuman(v.L*v.Lexpo+v.R*v.Rexpo+v.G*v.Gexpo+v.B*v.Bexpo+v.S*v.Sexpo+v.H*v.Hexpo+v.O*v.Oexpo))
			if v.L > 0 {
				fmt.Printf("L\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.L, v.Lexpo, secondsToHuman(v.L*v.Lexpo))
			}
			if v.R > 0 {
				fmt.Printf("R\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.R, v.Rexpo, secondsToHuman(v.R*v.Rexpo))
			}
			if v.G > 0 {
				fmt.Printf("G\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.G, v.Gexpo, secondsToHuman(v.G*v.Gexpo))
			}
			if v.B > 0 {
				fmt.Printf("B\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.B, v.Bexpo, secondsToHuman(v.B*v.Bexpo))
			}
			if v.S > 0 {
				fmt.Printf("S\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.S, v.Sexpo, secondsToHuman(v.S*v.Sexpo))
			}
			if v.H > 0 {
				fmt.Printf("H\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.H, v.Hexpo, secondsToHuman(v.H*v.Hexpo))
			}
			if v.O > 0 {
				fmt.Printf("O\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.O, v.Oexpo, secondsToHuman(v.O*v.Oexpo))
			}
			fmt.Println()
		}
	}

	if wdest.writeToFile {

		dest := fmt.Sprintf("%s/Lights_Report.txt", *lightsdir)
		report, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer report.Close()
		buff := bufio.NewWriter(report)

		fmt.Fprintf(buff, "Targets list: %q\n\n", targets)
		for _, v := range *ts {
			fmt.Fprintf(buff, "Object: %-30s %s%s\n", v.name, "Total: ", secondsToHuman(v.L*v.Lexpo+v.R*v.Rexpo+v.G*v.Gexpo+v.B*v.Bexpo+v.S*v.Sexpo+v.H*v.Hexpo+v.O*v.Oexpo))
			if v.L > 0 {
				fmt.Fprintf(buff, "L\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.L, v.Lexpo, secondsToHuman(v.L*v.Lexpo))
			}
			if v.R > 0 {
				fmt.Fprintf(buff, "R\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.R, v.Rexpo, secondsToHuman(v.R*v.Rexpo))
			}
			if v.G > 0 {
				fmt.Fprintf(buff, "G\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.G, v.Gexpo, secondsToHuman(v.G*v.Gexpo))
			}
			if v.B > 0 {
				fmt.Fprintf(buff, "B\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.B, v.Bexpo, secondsToHuman(v.B*v.Bexpo))
			}
			if v.S > 0 {
				fmt.Fprintf(buff, "S\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.S, v.Sexpo, secondsToHuman(v.S*v.Sexpo))
			}
			if v.H > 0 {
				fmt.Fprintf(buff, "H\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.H, v.Hexpo, secondsToHuman(v.H*v.Hexpo))
			}
			if v.O > 0 {
				fmt.Fprintf(buff, "O\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.O, v.Oexpo, secondsToHuman(v.O*v.Oexpo))
			}
			fmt.Fprintln(buff)
		}
		buff.Flush()
		report.Sync()
		report.Close()
	}
}

func plural(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		//result = strconv.Itoa(count) + " " + singular + " "
		result = fmt.Sprintf("%02d %s  ", count, singular)
	} else {
		//result = strconv.Itoa(count) + " " + singular + "s "
		result = fmt.Sprintf("%02d %ss ", count, singular)
	}
	return
}

func secondsToHuman(input int) (result string) {
	hours := math.Floor(float64(input) / 60 / 60)
	seconds := input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60
	if hours > 0 {
		result = plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else if minutes > 0 {
		result = plural(int(minutes), "minute") + plural(int(seconds), "second")
	} else {
		result = plural(int(seconds), "second")
	}
	return
}
