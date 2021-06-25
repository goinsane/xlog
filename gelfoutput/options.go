package gelfoutput

// Options defines several GELF options.
type Options struct {
	Address  string
	UseTCP   bool
	Host     string
	Facility string
}
