package meta

// this data is normally populated during build using -ldflags -X ...
// see magefile.go
var (
	Version = "unknown-dev"
	GitHash = "unknown"
)
