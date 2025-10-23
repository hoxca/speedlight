package utils

import Log "github.com/apatters/go-conlog"

func SetUpLogs(verbosity string) {
	formatter := Log.NewStdFormatter()
	formatter.Options.LogLevelFmt = Log.LogLevelFormatLongTitle
	Log.SetFormatter(formatter)
	switch verbosity {
	case "debug":
		Log.SetLevel(Log.DebugLevel)
	case "info":
		Log.SetLevel(Log.InfoLevel)
	case "warn":
		Log.SetLevel(Log.WarnLevel)
	case "error":
		Log.SetLevel(Log.ErrorLevel)
	default:
		Log.SetLevel(Log.WarnLevel)
	}
}
