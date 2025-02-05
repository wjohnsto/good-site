package pages

import (
	"path"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"good.site/pkg/config"
	"good.site/pkg/services"
	"good.site/pkg/services/rss"
	"good.site/pkg/utils/convert"
)

func init() {
	rootPath := "/posts/"
	cfg := config.GetConfig()
	markdownPath := cfg.Build.Markdown
	files := services.GetMarkdownFiles(markdownPath+rootPath, func(url string, createdDate time.Time) string {
		return path.Join(rootPath, createdDate.Format(strings.ReplaceAll(time.DateOnly, "-", "/")), url, "/")
	})

	description := "Blog posts."
	posts := []MarkdownPost{}
	singleOptions := MarkdownPageOptions{
		PublishedTerm: "Published",
		ShowAuthor:    true,
		RootUrl:       rootPath,
		DateFormat:    "January 02, 2006",
	}
	listOptions := ListMarkdownPageOptions{
		RootUrl:          singleOptions.RootUrl,
		Title:            "Posts",
		Description:      description,
		DateFormat:       "January 02, 2006",
		LinkToSinglePage: true,
		RenderHtml:       false,
	}

	for _, file := range files {
		post := MarkdownPost{
			Title:       file.Data.Title,
			TitleSlug:   slug.Make(file.Data.Title),
			Description: file.Data.Description,
			Slug:        file.Data.Slug,
			Url:         convert.ToSafeUrl(file.Data.Slug, true),
			Author:      file.Data.Author,
			CreatedDate: file.Data.CreatedDate,
			UpdatedDate: file.Data.UpdatedDate,
			Html:        file.Html,
			Text:        file.Text,
			Content:     convert.ToUnsafe(file.Html),
		}
		AddRoute(post.Slug, func() Page {
			return NewMarkdownPage(&post, &singleOptions)
		})
		posts = append(posts, post)
	}

	AddRoute(listOptions.RootUrl, func() Page {
		return NewListMarkdownPage(posts, &listOptions)
	})

	rss.NewFeedGenerator(&rss.FeedGeneratorOptions{
		RootPath:         cfg.Build.Public,
		RelativeFeedPath: rootPath,
		SiteUrl:          cfg.Site.Url,
		Title:            cfg.Site.Name + " blog",
		Description:      description,
		Author:           cfg.Site.Author,
		Language:         cfg.Site.Language,
		Copyright:        cfg.Site.Copyright,
		Files:            files,
		IconSm:           cfg.Site.Icons.Sm,
		IconMd:           cfg.Site.Icons.Md,
		IconLg:           cfg.Site.Icons.Lg,
		Logo2x1:          cfg.Site.Logo2x1,
	}).Generate()
}
