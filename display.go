package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

type writeDestination struct {
	writeToConsole bool
	writeToFile    bool
}

func (o *Object) printObject() {
	if o.rotation == 666 {
		fmt.Printf("Object name: %-38s %-9s Rotation: N/A\n", o.name, " ")
	} else {
		fmt.Printf("Object name: %-38s %-10s Rotation:%-3d°\n", o.name, " ", o.rotation)
	}
	ts := o.targets
	for _, v := range ts {
		fmt.Println()
		fmt.Printf("%d°C%-34s Total:%s",
			v.temp,
			" ",
			secondsToHuman((v.fltr.L+v.fltr.R+v.fltr.G+v.fltr.B+v.fltr.S+v.fltr.H+v.fltr.O)*v.expo))
		if v.fltr.L > 0 || v.fltr.R > 0 || v.fltr.G > 0 || v.fltr.B > 0 {
			fmt.Println()
		}
		if v.fltr.L > 0 {
			fmt.Printf("L\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.L, v.expo, secondsToHuman(v.fltr.L*v.expo))
		}
		if v.fltr.R > 0 {
			fmt.Printf("R\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.R, v.expo, secondsToHuman(v.fltr.R*v.expo))
		}
		if v.fltr.G > 0 {
			fmt.Printf("G\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.G, v.expo, secondsToHuman(v.fltr.G*v.expo))
		}
		if v.fltr.B > 0 {
			fmt.Printf("B\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.B, v.expo, secondsToHuman(v.fltr.B*v.expo))
		}
		if v.fltr.S > 0 || v.fltr.H > 0 || v.fltr.O > 0 {
			fmt.Println()
		}
		if v.fltr.S > 0 {
			fmt.Printf("S\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.S, v.expo, secondsToHuman(v.fltr.S*v.expo))
		}
		if v.fltr.H > 0 {
			fmt.Printf("H\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.H, v.expo, secondsToHuman(v.fltr.H*v.expo))
		}
		if v.fltr.O > 0 {
			fmt.Printf("O\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.O, v.expo, secondsToHuman(v.fltr.O*v.expo))
		}
	}
}

func (o *Object) fprintObject(buff io.Writer) {
	if o.rotation == 666 {
		fmt.Fprintf(buff, "Object name: %-38s %-9s Rotation: N/A\n", o.name, " ")
	} else {
		fmt.Fprintf(buff, "Object name: %-38s %-10s Rotation:%-3d°\n", o.name, " ", o.rotation)
	}
	ts := o.targets
	for _, v := range ts {
		fmt.Fprintln(buff)
		fmt.Fprintf(buff, "%d°C%-34s Total:%s",
			v.temp,
			" ",
			secondsToHuman((v.fltr.L+v.fltr.R+v.fltr.G+v.fltr.B+v.fltr.S+v.fltr.H+v.fltr.O)*v.expo))
		if v.fltr.L > 0 || v.fltr.R > 0 || v.fltr.G > 0 || v.fltr.B > 0 {
			fmt.Fprintln(buff)
		}
		if v.fltr.L > 0 {
			fmt.Fprintf(buff, "L\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.L, v.expo, secondsToHuman(v.fltr.L*v.expo))
		}
		if v.fltr.R > 0 {
			fmt.Fprintf(buff, "R\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.R, v.expo, secondsToHuman(v.fltr.R*v.expo))
		}
		if v.fltr.G > 0 {
			fmt.Fprintf(buff, "G\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.G, v.expo, secondsToHuman(v.fltr.G*v.expo))
		}
		if v.fltr.B > 0 {
			fmt.Fprintf(buff, "B\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.B, v.expo, secondsToHuman(v.fltr.B*v.expo))
		}
		if v.fltr.S > 0 || v.fltr.H > 0 || v.fltr.O > 0 {
			fmt.Fprintln(buff)
		}
		if v.fltr.S > 0 {
			fmt.Fprintf(buff, "S\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.S, v.expo, secondsToHuman(v.fltr.S*v.expo))
		}
		if v.fltr.H > 0 {
			fmt.Fprintf(buff, "H\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.H, v.expo, secondsToHuman(v.fltr.H*v.expo))
		}
		if v.fltr.O > 0 {
			fmt.Fprintf(buff, "O\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.fltr.O, v.expo, secondsToHuman(v.fltr.O*v.expo))
		}
	}
}

func (obs *Objects) printObjects(wdest writeDestination) {
	objects := obs.getObjects()

	if wdest.writeToConsole {
		fmt.Printf("Targets list: %q\n\n", objects)
		for _, v := range *obs {
			fmt.Println()
			v.printObject()
		}
		fmt.Println()
	}

	if wdest.writeToFile {
		dest := fmt.Sprintf("%s/Lights_Report.txt", *lightsdir)
		report, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			logFatal(err)
		}
		defer report.Close()

		buff := bufio.NewWriter(report)
		fmt.Fprintf(buff, "Targets list: %q\n\n", objects)

		for _, v := range *obs {
			fmt.Fprintln(buff)
			v.fprintObject(buff)

			fmt.Fprintln(buff)
		}
		buff.Flush()
		_ = report.Sync()
		report.Close()
	}
}

func plural(count int, singular string) string {
	var result string
	if (count == 1) || (count == 0) {
		result = fmt.Sprintf("%02d %s  ", count, singular)
	} else {
		result = fmt.Sprintf("%02d %ss ", count, singular)
	}
	return result
}

func secondsToHuman(input int) string {
	var result string
	hours := math.Floor(float64(input) / 60 / 60)
	seconds := input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60
	if hours > 0 {
		result = plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(seconds, "second")
	} else if minutes > 0 {
		result = plural(int(hours), "hour") + plural(int(minutes), "minute") + plural(seconds, "second")
	} else {
		result = plural(seconds, "second")
	}
	return result
}
