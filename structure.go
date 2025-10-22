package main

type target struct {
	name   string
	L      int
	Lexpo  int
	R      int
	Rexpo  int
	G      int
	Gexpo  int
	B      int
	Bexpo  int
	R1     int
	R1expo int
	G1     int
	G1expo int
	B1     int
	B1expo int
	S      int
	Sexpo  int
	H      int
	Hexpo  int
	O      int
	Oexpo  int
}

type targets map[string]target

var targetList = targets{}
var o *target

func (ts targets) set(key string, value *target) {
	ts[key] = *value
}

/*
  func (ts targets) get(key string) *target {
          obj := ts[key]
          return &obj
  }
*/

func (ts targets) exist(key string) bool {
	for k := range ts {
		if k == key {
			return true
		}
	}
	return false
}

func newTarget(object string) *target {
	return &target{
		name: object,
		L:    0,
		R:    0,
		G:    0,
		B:    0,
		R1:   0,
		G1:   0,
		B1:   0,
		S:    0,
		H:    0,
		O:    0,
	}
}

func (ts *targets) getTargets() []string {
	ret := []string{}
	for _, v := range *ts {
		ret = append(ret, v.name)
	}
	return ret
}

func (t *target) iterateFilter(filter string, expo int) {
	switch filter {
	case "L":
		t.L++
		t.Lexpo = expo
	case "R":
		if expo > 240 {
			t.R++
			t.Rexpo = expo
		} else {
			t.R1++
			t.R1expo = expo
		}
	case "G":
		if expo > 240 {
			t.G++
			t.Gexpo = expo
		} else {
			t.G1++
			t.G1expo = expo
		}
	case "B":
		if expo > 240 {
			t.B++
			t.Bexpo = expo
		} else {
			t.B1++
			t.B1expo = expo
		}
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
