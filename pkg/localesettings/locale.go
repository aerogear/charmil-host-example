package localesettings

import "embed"

// DefaultLocales stores the embedded contents of all the locales files
//go:embed locales/*
var DefaultLocales embed.FS
