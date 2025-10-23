# Speedlight

Speedligh give a flash view on your astronomical target lights.

Speedlight is a simple little binary to maintain counts and statitiscs of lights in your acquisition data directory.
You just need to provide the path to your incoming Voyager lights (-dir path).

Speedlight should evolve to a simple grader to avoid counting the bad lights in the reports.

# Usage

Usage of ./bin/speedlight:

  -dir string       lights directory (default "D:/Data/Voyager/Lights/")
  -level string     set log level of speedlight default warn (default "warn")
  -report           write report to the filesystem (default true)
  -rotation         tell if rotation is used (default true)

