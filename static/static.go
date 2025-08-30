package static

import (
	"embed"
	"io/fs"
)

//go:embed *
var embeddedFiles embed.FS

var Files, _ = fs.Sub(embeddedFiles, ".") // Optional, to clean base path
