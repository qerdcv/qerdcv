package web

import "embed"

//go:embed template/*.html
var Template embed.FS

//go:embed assets/*/**
var Static embed.FS
