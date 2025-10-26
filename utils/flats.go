package utils

type flats map[float32]Filters

var FlatList = flats{}

var rotation float32
var Rotations []float32

var flt *Filters

var TimeFrame int

func addFlats(rotation float32, filter string) {
	if !FlatList.exist(rotation) {
		flt = newFlat()
		Rotations = appendIfMissing(Rotations, rotation)
		FlatList.set(rotation, flt)
	}
	flt.iterateFilter(filter)
	FlatList.set(rotation, flt)
}

func (fs flats) set(key float32, value *Filters) {
	fs[key] = *value
}

func appendIfMissing(slice []float32, i float32) []float32 {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func (f flats) exist(key float32) bool {
	for k := range f {
		if k == key {
			return true
		}
	}
	return false
}

func newFlat() *Filters {
	return &Filters{
		L: 0,
		R: 0,
		G: 0,
		B: 0,
		S: 0,
		H: 0,
		O: 0,
	}
}

func (f *Filters) iterateFilter(filter string) {
	switch filter {
	case "L":
		f.L++
	case "R":
		f.R++
	case "G":
		f.G++
	case "B":
		f.B++
	case "S":
		f.S++
	case "H":
		f.H++
	case "O":
		f.O++
	}
}

func (f Filters) String() (filters string) {
	filters = ""

	if f.L > 0 {
		filters = filters + "L"
	}
	if f.R > 0 {
		filters = filters + "R"
	}
	if f.G > 0 {
		filters = filters + "G"
	}
	if f.B > 0 {
		filters = filters + "B"
	}
	if f.S > 0 {
		filters = filters + "S"
	}
	if f.H > 0 {
		filters = filters + "H"
	}
	if f.O > 0 {
		filters = filters + "O"
	}

	return filters
}
