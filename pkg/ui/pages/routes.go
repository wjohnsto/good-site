package pages

import (
	"log/slog"
	"strings"

	"github.com/a-h/templ"
)

type Page interface {
	Render() templ.Component
}

type NewPageFunction func() Page

type Route struct {
	Path string
	New  NewPageFunction
}

var Routes = []*Route{}

func AddRoute(path string, new func() Page) {
	slog.Debug("Adding route: " + path)

	Routes = append(Routes, &Route{
		Path: strings.Trim(path, "/") + "/",
		New:  new,
	})
}
