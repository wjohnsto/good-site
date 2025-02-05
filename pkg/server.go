package server

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"good.site/pkg/config"
	"good.site/pkg/ui/pages"
	"good.site/pkg/utils/convert"
	"good.site/pkg/utils/file"
)

type Server struct {
	publicPath string
}

func (s *Server) RenderRoute(ctx context.Context, route *pages.Route) {
	dir := path.Join(s.publicPath, route.Path)
	name := dir

	if filepath.Ext(name) == "" {
		name = path.Join(dir, "index.html")
	}

	buf := &bytes.Buffer{}
	t := route.New().Render()

	if err := t.Render(ctx, buf); err != nil {
		panic(err)
	}
	slog.Debug(fmt.Sprintf("%s created", name))

	doc, err := convert.ParseHtml(buf.String())

	if err != nil {
		panic(err)
	}

	css := convert.ExtractStyles(doc)

	file.WriteFile(name, convert.RenderNode(doc))
	file.WriteFile(path.Join(name, "../styles.css"), strings.TrimSpace(css))
}

func (s *Server) Render() {
	ctx := context.Background()

	if err := os.MkdirAll(s.publicPath, os.ModePerm); err != nil && err != os.ErrExist {
		panic(err)
	}

	logs := []string{}

	for _, route := range pages.Routes {
		if strings.Contains(route.Path, "[") {
			continue
		}

		s.RenderRoute(ctx, route)
		logs = append(logs, route.Path)
	}

	slog.Info(fmt.Sprintf("Creating %d pages:\n  %s", len(logs), strings.Join(logs, "\n  ")))
}

func New() *Server {
	server := &Server{
		publicPath: config.GetConfig().Build.Public,
	}

	return server
}
