package demo_app

import "runtime/debug"

func init() {
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" {
				GitCommit = s.Value
				break
			}
		}
	}
}

var (
	Version   = "0.0.1"
	GitCommit = "unknown"
)
