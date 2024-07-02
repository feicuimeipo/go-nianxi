package web

import (
	"embed"
)

//go:embed dist
var StaticFs embed.FS

////go:embed dist
//var TemplateFs embed.FS
