package utils

var Flags struct {
	WithinOpenWhisk bool   // is this running within an OpenWhisk action?
	ApiHost         string // OpenWhisk API host
	Auth            string // OpenWhisk API key
	ApiVersion      string // OpenWhisk version
}
