package version

var (
	// Version is the version number, set by ldflags during build
	Version = "dev"
	// GitCommit is the git commit hash, set by ldflags during build
	GitCommit = "unknown"
	// BuildTime is the build timestamp, set by ldflags during build
	BuildTime = "unknown"
)
