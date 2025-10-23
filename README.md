# Speedlight

Speedligh give a flash view on your astronomical target lights.

Speedlight is a simple little binary to maintain counts and statitiscs of lights in your acquisition data directory.
You just need to provide the path to your incoming Voyager lights (-dir path).

Speedlight should evolve to a simple grader to avoid counting the bad lights in the reports.

# Usage

Usage of ./bin/speedlight:

&nbsp;&nbsp;&nbsp;-dir&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;string&nbsp;&nbsp;lights directory (default "D:/Data/Voyager/Lights/") <br/>
&nbsp;&nbsp;&nbsp;-level&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;string&nbsp;  set log level of speedlight default warn (default "warn")<br/>
&nbsp;&nbsp;&nbsp;-report&nbsp;&nbsp;&nbsp;&nbsp;bool&nbsp;&nbsp; write report to the filesystem (default true)<br/>
&nbsp;&nbsp;&nbsp;-rotation bool&nbsp;&nbsp; tell if rotation is used (default true)<br/>

