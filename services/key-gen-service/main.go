package main

import (
	"context"
	"os"

	versioninfo "github.com/sazonovItas/go-pastebin/pkg/version"
	cmd "github.com/sazonovItas/go-pastebin/services/key-gen-service/cmd/keygen"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	context.Background()
	vInfo := versioninfo.BuildVersion(
		"Key gen service",
		"Key gen service generates unique keys",
		version,
		commit,
		date,
	)
	cmd.Execute(vInfo, os.Args[1:])
}
