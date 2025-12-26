package version

// These values are set via -ldflags at build time by GoReleaser.
// Defaults are useful for local builds.
var (
	// semver from tag, e.g. v0.1.0
	Version = "v0.0.0-dev"
	// SHA of the build
	Commit = "none"
	// RFC3339 date of the build
	Date = "unknown"
)
