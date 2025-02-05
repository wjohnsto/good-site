package main

import (
	"path"
	"regexp"

	"good.site/pkg/utils/convert"
	"good.site/pkg/utils/file"
)

func main() {
	r, err := regexp.Compile(`index\.html$`)

	if err != nil {
		panic(err)
	}

	files := file.FindFiles(convert.ToPath("public"), r)

	for _, p := range files {
		cssPath := path.Join(p, "../styles.css")
		htmlStr := file.ReadFile(p)
		cssStr := file.ReadFile(cssPath)

		doc, err := convert.ParseHtml(string(htmlStr[:]))

		if err != nil {
			panic(err)
		}

		convert.AddStyle(doc, string(cssStr[:]))

		outStr := convert.RenderNode(doc)

		file.WriteFile(p, outStr)
		err = file.DeleteFile(cssPath)

		if err != nil {
			panic(err)
		}
	}
}
