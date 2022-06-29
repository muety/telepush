package views

import (
	"embed"
	"github.com/muety/telepush/config"
	"io/fs"
	"os"
)

//go:embed static
var staticFiles embed.FS

//go:embed *.html
var templateFiles embed.FS

func GetStaticFilesFS() (fsys fs.FS) {
	fsys, _ = fs.Sub(staticFiles, "static")
	if config.Get().IsDev() {
		fsys = os.DirFS("views/static")
	}
	return fsys
}

func GetTemplatesFS() (fsys fs.FS) {
	fsys = templateFiles
	if config.Get().IsDev() {
		fsys = os.DirFS("views")
	}
	return fsys
}
