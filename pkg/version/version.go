package version

import goversion "github.com/caarlos0/go-version"

func BuildVersion(name, desc, version, commit, date string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails(name, desc, ""),
		func(i *goversion.Info) {
			if version != "" {
				i.GitVersion = version
			}
			if commit != "" {
				i.GitCommit = commit
			}
			if date != "" {
				i.BuildDate = date
			}
		},
	)
}
