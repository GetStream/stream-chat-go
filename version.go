package stream_chat

import (
	"fmt"
)

const (
	versionMajor = 6
	versionMinor = 7
	versionPatch = 8
)

// Version returns the version of the library.
func Version() string {
	return "v" + fmtVersion()
}

func versionHeader() string {
	return "stream-go-client-" + fmtVersion()
}

func fmtVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		versionMajor,
		versionMinor,
		versionPatch)
}
