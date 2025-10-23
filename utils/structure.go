package utils

import (
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
	rotation int
	targets  map[string]Target
}

type Objects map[string]Object
type Targets map[string]Target

var ObjectList = Objects{}
var targetList = Targets{}
var t *Target
var o *Object

func addTarget(target string, rotation int) *Object {
	targetList = Targets{}
	result := strings.Split(target, "~")
	targetName := result[0]
	if !ObjectList.exist(targetName) {
		Log.Debugf("create object %s", targetName)
		o = newObject(target, rotation)
		ObjectList.set(targetName, o)
		// Log.Debugf("target Object: %q\n", o) .
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

func (os Objects) set(key string, value *Object) {
	os[key] = *value
}

func (o Object) setTarget(key string, value *Target) {
	o.targets[key] = *value
}

func (os Objects) exist(key string) bool {
	for k := range os {
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

func (os *Objects) getObjects() []string {
	ret := []string{}
	for _, v := range *os {
		ret = append(ret, v.name)
	}
	sort.Strings(ret)
	return ret
}

func (ts *Targets) getTarget(target string) *Target {
	for _, v := range *ts {
		if v.tuple == target {
			return &v
		}
	}
	return nil
}

func (os *Objects) getObject(object string) *Object {
	// Log.Debugf("Search for object: %s", object) .
	for _, v := range *os {
		// Log.Debugf("get target: %s from %s\n", object, v.name) .
		if v.name == object {
			return &v
		}
	}
	return nil
}

// Constructors .

func newObject(target string, targetRotation int) *Object {
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
