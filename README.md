# Speedlight

Speedligh give a flash view on your astronomical target lights.

Speedlight is a simple little binary to maintain counts and statitiscs of lights in your acquisition data directory.
You just need to provide the path to your incoming Voyager lights (-dir path).

Speedlight should evolve to a simple grader to avoid counting the bad lights in the reports.

# Config file

config file is by default place in the `../conf/speedlight.yaml` of the program directory

```
lightsdir: '/Volumes/Dyno/Kosmodrom/Lights'
regexp: '/(.*)/[[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}/(.*)/.*_LIGHT_[LRGBSHO]*_([[:digit:]]*)s_BIN1_(.*)C_GA.*_[[:digit:]]{8}_[[:digit:]]{6}_[[:digit:]]{3}_PA([[:digit:]]{3}\.[[:digit:]]{2})_[EW]\.FIT'
time_frame: 14
level: "warn"
```

 * the `lightsdir` is the root directory of your lights
 * `time_frame` is the number of hour used for scanning lights files
    detect all files acquired in the last 14 hours
 * regexp is the regular expression used to capture the revelant informations
   you can use any regexp but the field capture order is important !
 * log level

```
    field 1: target object name
    field 2: filter used [LRGBSHO]
    field 3: exposure time
    field 4: camera temperature
    field 5: rotation angle
```

# Usage

Usage of ./bin/speedlight:

## report command 

```
hugh⨕shupo:speedlight [ main | ✚ 6 ♟ 2 ]|● ./bin/speedlight report -h
report will generate a report on all the lights
produced by your voyager astronomy orchestrator

it will sum time exposure by target and temperature.

Usage:
  speedlight report [flags]

Flags:
      --console    write report to the console (default true)
  -h, --help       help for report
      --report     write report to the filesystem (default true)
      --rotation   manage rotation in lights report (default true)

Global Flags:
      --config string   config file (default is conf/speedlight.yaml)
      --dir string      lights directory
      --level string    set log level
```

## rotation command

```
hugh⨕shupo:speedlight [ main | ✚ 6 ♟ 1 ]|● ./bin/speedlight rotation -h
Get the rotation used by the target during last acquisition night

Usage:
  speedlight rotation [flags]

Flags:
  -h, --help         help for rotation
      --target int   night target number, between 1 and 3

Global Flags:
      --config string   config file (default is conf/speedlight.yaml)
      --dir string      lights directory
      --level string    set log level
```

## filters command

```
hugh⨕shupo:speedlight [ main | ✚ 6 ♟ 1 ]|● ./bin/speedlight filters -h
Get the list of filters used for this target during the last acquisition night

note that if --target=0 is used, this command return all filters used during
the last acquisition night (all targets and rotations)
 
Usage:
  speedlight filters [flags]

Flags:
  -h, --help         help for filters
      --target int   night target number, between 0 and 3

Global Flags:
      --config string   config file (default is conf/speedlight.yaml)
      --dir string      lights directory
      --level string    set log level
```
