package demo_app

import "runtime/debug"

var (
	Version   = "0.0.2"
	GitCommit = "unknown2"
)

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
