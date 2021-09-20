package xlog

// Flag holds single or multiple flogs of Log.
type Flag int

const (
	FlagDate         = Flag(1 << iota)                                                                // prints the date in the local time zone: 2009/01/23
	FlagTime                                                                                          // prints the time in the local time zone: 01:23:23
	FlagMicroseconds                                                                                  // prints microsecond resolution: 01:23:23.123123
	FlagUTC                                                                                           // uses UTC rather than the local time zone
	FlagSeverity                                                                                      // prints severity level
	FlagPadding                                                                                       // prints padding with multiple lines
	FlagLongFunc                                                                                      // prints full package name and function name: a/b/c/d.Func1()
	FlagShortFunc                                                                                     // prints final package name and function name: d.Func1()
	FlagLongFile                                                                                      // prints full file name and line number: a/b/c/d.go:23
	FlagShortFile                                                                                     // prints final file name element and line number: d.go:23
	FlagFields                                                                                        // prints fields
	FlagStackTrace                                                                                    // prints stack trace
	FlagDefault      = FlagDate | FlagTime | FlagSeverity | FlagPadding | FlagFields | FlagStackTrace // holds initial flags for the Logger
)
