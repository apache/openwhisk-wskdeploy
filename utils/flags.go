package utils

var Flags struct {
	WithinOpenWhisk bool   // is this running within an OpenWhisk action?
	ApiHost         string // OpenWhisk API host
	Auth            string // OpenWhisk API key
	ApiVersion      string // OpenWhisk version

	//action flag definition
	//from go cli
	action struct {
		docker   bool
		copy     bool
		pipe     bool
		web      string
		sequence bool
		timeout  int
		memory   int
		logsize  int
		result   bool
		kind     string
		main     string
	}
}
