package version

import "fmt"

var (
	BuildDate    = "<unknown>"
	BuildVersion = "<unknown>"
)

func GetBuildDate() string {
	return BuildDate
}

func GetBuildVersion() string {
	return BuildVersion
}

func GetVersionString() string {
	return fmt.Sprintf("Built on %s from %s", BuildDate, BuildVersion)
}
