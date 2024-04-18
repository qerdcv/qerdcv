package web

import "embed"

var (

	//go:embed templates
	Templates embed.FS

	//go:embed static
	Static embed.FS
)
