package main

import "fmt"

/*
type writeDestination struct {
	writeToConsole bool
	writeToFile    bool
}
*/

func (os *Objects) printObjects() {
	objects := os.getObjects()
	fmt.Printf("Targets list: %q\n", objects)
	for _, v := range *os {
		v.printObject()
	}
}

func (o *Object) printObject() {
	fmt.Printf("Object name: %s\n", o.name)
	ts := o.targets
	for _, v := range ts {
		v.printTarget()
	}
}

func (t *Target) printTarget() {
	fmt.Printf("\tTarget Name: %s\n\tTarget temperature: %d\n\tTarget exposure: %d\n\tTarget rotation: %d\n",
		t.name, t.temp, t.expo, t.rot)
	fmt.Printf("\tFilters: {\n \t L:%d\n \t R:%d\n \t G:%d\n \t B:%d\n \t S:%d\n \t H:%d\n \t O:%d\n\t}\n",
		t.fltr.L, t.fltr.R, t.fltr.G, t.fltr.B, t.fltr.S, t.fltr.H, t.fltr.O)
}

/*
func (ts *targets) printObjects(wdest writeDestination) {
	targets := ts.getTargets()

	if wdest.writeToConsole {
		fmt.Printf("Targets list: %q\n\n", targets)
		for _, v := range *ts {
			fmt.Println()
			fmt.Printf("Object: %-30s Total:%s\n",
				v.name,
				secondsToHuman(v.L*v.Lexpo+v.R*v.Rexpo+v.G*v.Gexpo+v.B*v.Bexpo+v.S*v.Sexpo+v.H*v.Hexpo+v.O*v.Oexpo))
			if v.L > 0 || v.R > 0 || v.G > 0 || v.B > 0 {
				fmt.Println()
			}
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
			if v.R1 > 0 || v.G1 > 0 || v.B1 > 0 {
				fmt.Println()
			}
			if v.R1 > 0 {
				fmt.Printf("R\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.R1, v.R1expo, secondsToHuman(v.R1*v.R1expo))
			}
			if v.G1 > 0 {
				fmt.Printf("G\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.G1, v.G1expo, secondsToHuman(v.G1*v.G1expo))
			}
			if v.B1 > 0 {
				fmt.Printf("B\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.B1, v.B1expo, secondsToHuman(v.B1*v.B1expo))
			}
			if v.S > 0 || v.H > 0 || v.O > 0 {
				fmt.Println()
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
			logFatal(err)
		}
		defer report.Close()

		buff := bufio.NewWriter(report)

		fmt.Fprintf(buff, "Targets list: %q\n\n", targets)
		for _, v := range *ts {
			fmt.Fprintln(buff)
			fmt.Fprintf(buff, "Object: %-30s Total: %s\n",
				v.name,
				secondsToHuman(v.L*v.Lexpo+v.R*v.Rexpo+v.G*v.Gexpo+v.B*v.Bexpo+v.S*v.Sexpo+v.H*v.Hexpo+v.O*v.Oexpo))
			if v.L > 0 || v.R > 0 || v.G > 0 || v.B > 0 {
				fmt.Fprintln(buff)
			}
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
			if v.R1 > 0 || v.G1 > 0 || v.B1 > 0 {
				fmt.Fprintln(buff)
			}
			if v.R1 > 0 {
				fmt.Fprintf(buff, "R\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.R1, v.R1expo, secondsToHuman(v.R1*v.R1expo))
			}
			if v.G1 > 0 {
				fmt.Fprintf(buff, "G\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.G1, v.G1expo, secondsToHuman(v.G1*v.G1expo))
			}
			if v.B1 > 0 {
				fmt.Fprintf(buff, "B\tNb: %4d\tExpo: %4ds\tSubs: %s\n", v.B1, v.B1expo, secondsToHuman(v.B1*v.B1expo))
			}
			if v.S > 0 || v.H > 0 || v.O > 0 {
				fmt.Fprintln(buff)
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
		result = plural(int(minutes), "minute") + plural(seconds, "second")
	} else {
		result = plural(seconds, "second")
	}
	return result
}
*/
