package xlog

// Flag holds single or multiple flags of Log.
// An Output instance uses these flags which are stored by Flag type.
type Flag int

const (
	// FlagDate prints the date in the local time zone: 2009/01/23
	FlagDate Flag = 1 << iota

	// FlagTime prints the time in the local time zone: 01:23:23
	FlagTime

	// FlagMicroseconds prints microsecond resolution: 01:23:23.123123
	FlagMicroseconds

	// FlagUTC uses UTC rather than the local time zone
	FlagUTC

	// FlagSeverity prints severity level
	FlagSeverity

	// FlagPadding prints padding with multiple lines
	FlagPadding

	// FlagLongFunc prints full package name and function name: a/b/c/d.Func1()
	FlagLongFunc

	// FlagShortFunc prints final package name and function name: d.Func1()
	FlagShortFunc

	// FlagLongFile prints full file name and line number: a/b/c/d.go:23
	FlagLongFile

	// FlagShortFile prints final file name element and line number: d.go:23
	FlagShortFile

	// FlagFields prints fields if there are
	FlagFields

	// FlagStackTrace prints stack trace if there is
	FlagStackTrace

	// FlagErfStackTrace prints stack traces of erf error if there are
	FlagErfStackTrace

	// FlagDefault holds initial flags for the Logger
	FlagDefault = FlagDate | FlagTime | FlagSeverity | FlagPadding | FlagFields | FlagStackTrace | FlagErfStackTrace
)
