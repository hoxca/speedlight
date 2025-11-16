package utils

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	Log "github.com/apatters/go-conlog"
)

type Filters struct {
	L int
	R int
	G int
	B int
	S int
	H int
	O int
}

type Target struct {
	tuple string
	name  string
	temp  int
	expo  int
	fltr  Filters
}

type Object struct {
	name     string
	rotation float32
	targets  map[string]Target
	SortBy   string
}

type Objects map[string]Object
type Targets map[string]Target

type by func(t1, t2 *Target) bool

var RotUsed bool

// compareAstronomicalNames compares astronomical object names with proper numeric sorting
// for postfix numbers (e.g., M42 comes after M5, not before)
func compareAstronomicalNames(a, b string) bool {
	// Extract prefix and numeric parts - more flexible pattern
	re := regexp.MustCompile(`^([A-Za-z]+)(\d*)$`)
	aMatch := re.FindStringSubmatch(a)
	bMatch := re.FindStringSubmatch(b)

	// If both match the pattern, compare by prefix then numerically
	if len(aMatch) == 3 && len(bMatch) == 3 {
		// Compare prefixes first
		if aMatch[1] != bMatch[1] {
			return aMatch[1] < bMatch[1]
		}

		// Compare numeric parts
		aNum := aMatch[2]
		bNum := bMatch[2]

		// Handle empty numeric parts
		if aNum == "" && bNum == "" {
			return a < b
		}
		if aNum == "" {
			return true // no number comes before number
		}
		if bNum == "" {
			return false // number comes after no number
		}

		// Convert to integers and compare
		aInt, errA := strconv.Atoi(aNum)
		bInt, errB := strconv.Atoi(bNum)
		if errA == nil && errB == nil {
			return aInt < bInt
		}
	}

	// Fallback to string comparison
	return a < b
}

var ObjectList = Objects{}
var targetList = Targets{}
var t *Target
var o *Object

func addTarget(target string, rotation float32) *Object {
	targetList = Targets{}
	result := strings.Split(target, "~")
	targetName := result[0]
	if !ObjectList.exist(targetName) {
		Log.Debugf("create object %s", targetName)
		o = newObject(target, rotation)
		ObjectList.set(targetName, o)
		// Log.Debugf("target Object: %q ", o) .
	}
	o = ObjectList.getObject(targetName)
	targetList = o.targets
	// Log.Debugf("targetList %q for Object Name: %s", targetList, o.name) .

	if !targetList.exist(target) {
		Log.Debugf("create target: %s", target)
		t = newTarget(target)
		o.setTarget(target, t)
		return o
	}
	t = targetList.getTarget(target)
	if t != nil {
		return o
	}
	return nil
}

func (o *Object) addFilter(target string, filter string) {
	targetList = Targets{}
	targetList = o.targets
	t = targetList.getTarget(target)

	t.iterateFilter(filter)
	o.setTarget(t.tuple, t)
}

func (t *Target) iterateFilter(filter string) {
	switch filter {
	case "L":
		t.fltr.L++
	case "R":
		t.fltr.R++
	case "G":
		t.fltr.G++
	case "B":
		t.fltr.B++
	case "S":
		t.fltr.S++
	case "H":
		t.fltr.H++
	case "O":
		t.fltr.O++
	}
}

// Getters Setters .

func (obs Objects) set(key string, value *Object) {
	obs[key] = *value
}

func (o Object) setTarget(key string, value *Target) {
	o.targets[key] = *value
}

func (obs Objects) exist(key string) bool {
	for k := range obs {
		if k == key {
			return true
		}
	}
	return false
}

func (ts Targets) exist(key string) bool {
	for k := range ts {
		if k == key {
			return true
		}
	}
	return false
}

func (obs *Objects) getObjects() []string {
	ret := []string{}
	for _, v := range *obs {
		ret = append(ret, v.name)
	}
	sort.Slice(ret, func(i, j int) bool {
		return compareAstronomicalNames(ret[i], ret[j])
	})
	return ret
}

func (o *Object) getTargets() []Target {
	var targets []Target
	for _, target := range o.targets {
		targets = append(targets, target)
	}
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].tuple < targets[j].tuple
	})
	return targets
}

func (ts *Targets) getTarget(target string) *Target {
	for _, v := range *ts {
		if v.tuple == target {
			return &v
		}
	}
	return nil
}

func (obs *Objects) getObject(object string) *Object {
	// Log.Debugf("Search for object: %s", object) .
	for k, v := range *obs {
		// Log.Debugf("get target: %s from %s ", object, v.name) .
		if v.name == object {
			v.SortBy = k
			return &v
		}
	}
	return nil
}

// Constructors .

func newObject(target string, targetRotation float32) *Object {
	return &Object{
		name:     strings.Split(target, "~")[0],
		rotation: targetRotation,
		targets:  targetList,
	}
}

func newTarget(target string) *Target {
	result := strings.Split(target, "~")
	targetName := result[0]
	targetTemperature, _ := strconv.Atoi(result[1])
	targetExposition, _ := strconv.Atoi(result[2])
	return &Target{
		tuple: target,
		name:  targetName,
		temp:  targetTemperature,
		expo:  targetExposition,
		fltr: Filters{
			L: 0,
			R: 0,
			G: 0,
			B: 0,
			S: 0,
			H: 0,
			O: 0,
		},
	}
}
