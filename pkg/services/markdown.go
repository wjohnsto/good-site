package services

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/djherbis/times"
	"github.com/gosimple/slug"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	ghtml "github.com/yuin/goldmark/renderer/html"
	"good.site/pkg/utils/convert"
	"good.site/pkg/utils/file"
)

type MetaData struct {
	Title       string
	Description string
	Author      string
	Slug        string
	CreatedDate time.Time
	UpdatedDate time.Time
}

type MarkdownFile struct {
	Data     MetaData
	FilePath string
	Html     string
	Text     string
}

func ParseDate(date string) time.Time {
	dt, err := time.Parse(time.RFC3339, date)

	if err != nil {
		dt, err = time.Parse(time.DateOnly, date)

		if err != nil {
			panic(err)
		}
	}

	return dt
}

func SetMeta(filePath string, text string, mdFile *MarkdownFile) {

	meta := "---"

	if mdFile.Data.Title != "" {
		meta = fmt.Sprintf("%s\ntitle: \"%s\"\nslug: \"%s\"", meta, mdFile.Data.Title, slug.Make(mdFile.Data.Title))
	}

	if mdFile.Data.Author != "" {
		meta = fmt.Sprintf("%s\nauthor: \"%s\"", meta, mdFile.Data.Author)
	}

	if mdFile.Data.Description != "" {
		meta = fmt.Sprintf("%s\ndescription: \"%s\"", meta, mdFile.Data.Description)
	}

	meta = fmt.Sprintf(
		"%s\ncreated_date: \"%s\"\nupdated_date: \"%s\"\n---\n",
		meta,
		mdFile.Data.CreatedDate.Format(time.RFC3339),
		mdFile.Data.UpdatedDate.Format(time.RFC3339),
	)

	split := strings.Split(text, "\n")
	metaStart := -1
	metaEnd := -1

	for index, line := range split {
		if len(line) < 3 {
			continue
		}

		if line[0:3] == "---" {
			if metaStart < 0 {
				metaStart = index
			} else {
				metaEnd = index
				break
			}
		}
	}

	if metaStart > -1 && metaEnd > -1 {
		text = strings.Join(split[metaEnd+1:], "\n")
	}

	file.WriteFile(filePath, meta+text)
}

func ParseMarkdown(filePath string, getUrl func(slug string, createdDate time.Time) string) *MarkdownFile {
	dat, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	fi, err := times.Stat(filePath)

	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	context := parser.NewContext()

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, meta.Meta, highlighting.NewHighlighting(
			highlighting.WithStyle("base16-snazzy"),
			highlighting.WithFormatOptions(
				chromahtml.ClassPrefix("highlight "),
				chromahtml.WithClasses(true),
				chromahtml.LineNumbersInTable(true),
				chromahtml.WithAllClasses(true),
				chromahtml.WithLineNumbers(true),
			),
		)),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			ghtml.WithHardWraps(),
			ghtml.WithXHTML(),
		),
	)

	if err := md.Convert(dat, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	doc, err := convert.ParseHtml(buf.String())

	if err != nil {
		panic(err)
	}

	title := convert.TextContentFromHtml(doc, "h1")
	description := convert.TextContentFromHtml(doc, "p")

	data := meta.Get(context)

	if data["title"] != nil {
		title = data["title"].(string)
	}

	if len(title) == 0 {
		panic(fmt.Errorf("error with file: %s, no title found", filePath))
	}

	convert.RemoveNodes(doc, "h1", 1)
	contentHtml := convert.RenderNode(doc)
	url := slug.Make(title)

	if data["slug"] != nil {
		url = data["slug"].(string)
	}

	if data["description"] != nil {
		description = data["description"].(string)
	}

	updatedDate := fi.ModTime()
	createdDate := fi.BirthTime()

	if data["created_date"] != nil {
		createdDate = ParseDate(data["created_date"].(string))
	}

	if data["updated_date"] != nil {
		updatedDate = ParseDate(convert.DefaultTo(data["updated_date"], createdDate.Format(time.RFC3339)).(string))
	}

	file := MarkdownFile{
		Data: MetaData{
			Title:       title,
			Description: description,
			Slug:        getUrl(url, createdDate),
			Author:      convert.DefaultTo(data["author"], "").(string),
			CreatedDate: createdDate,
			UpdatedDate: updatedDate,
		},
		FilePath: filePath,
		Html:     contentHtml,
		Text:     string(dat),
	}

	SetMeta(filePath, string(dat[:]), &file)

	return &file
}

func GetMarkdownFiles(dir string, getUrl func(slug string, createdDate time.Time) string) []*MarkdownFile {
	files := []*MarkdownFile{}

	r, err := regexp.Compile(`\.md$`)

	if err != nil {
		panic(err)
	}

	paths := file.FindFiles(convert.ToPath(dir), r)

	for _, path := range paths {
		f := ParseMarkdown(path, getUrl)

		if f != nil {
			files = append(files, f)
		}
	}

	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Data.CreatedDate.After(files[j].Data.CreatedDate)
	})

	return files
}
