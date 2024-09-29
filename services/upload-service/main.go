package main

import (
	"os"

	versioninfo "github.com/sazonovItas/go-pastebin/pkg/version"
	cmd "github.com/sazonovItas/go-pastebin/services/upload-service/cmd/upload"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	vInfo := versioninfo.BuildVersion(
		"Upload service",
		"Upload service",
		version,
		commit,
		date,
	)
	cmd.Execute(vInfo, os.Args[1:])
}
