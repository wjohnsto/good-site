package rss

import (
	"encoding/json"
	"encoding/xml"
	"path"
	"regexp"
	"time"

	"good.site/pkg/services"
	"good.site/pkg/utils/convert"
	"good.site/pkg/utils/file"
)

type FeedGeneratorOptions struct {
	RootPath         string
	RelativeFeedPath string
	Title            string
	Description      string
	Author           string
	SiteUrl          string
	Language         string
	Copyright        string
	Files            []*services.MarkdownFile
	IconSm           string
	IconMd           string
	IconLg           string
	Logo2x1          string
}

type FeedGenerator struct {
	Options *FeedGeneratorOptions
	Atom    *AtomFeed
	Rss     *RssFeed
	Json    *JsonFeed
}

func (fg *FeedGenerator) getAtomItems() []AtomEntry {
	entries := []AtomEntry{}

	for _, file := range fg.Options.Files {
		entries = append(entries, AtomEntry{
			Title: file.Data.Title,
			Link: []AtomLink{{
				Rel:  "self",
				Href: convert.FormatLink(fg.Options.SiteUrl, file.Data.Slug),
			}},
			Id:        convert.FormatLink(fg.Options.SiteUrl, file.Data.Slug),
			Updated:   atomDate(file.Data.UpdatedDate),
			Published: atomDate(file.Data.CreatedDate),
			Summary:   file.Data.Description,
			Content: AtomContent{
				Type: "html",
				Html: file.Html,
			},
			Author: &AtomAuthor{
				Name: convert.DefaultTo(file.Data.Author, fg.Options.Author).(string),
			},
		})
	}

	return entries
}

func (fg *FeedGenerator) getRssItems() []RssItem {
	items := []RssItem{}

	for _, file := range fg.Options.Files {
		items = append(items, RssItem{
			Title: file.Data.Title,
			Link:  convert.FormatLink(fg.Options.SiteUrl, file.Data.Slug),
			Description: RssItemDescription{
				Content: file.Html,
				Type:    "html",
			},
			Author:  convert.DefaultTo(file.Data.Author, fg.Options.Author).(string),
			PubDate: rssDate(file.Data.CreatedDate),
			Guid:    convert.FormatLink(fg.Options.SiteUrl, file.Data.Slug),
		})
	}

	return items
}

func (fg *FeedGenerator) getJsonItems() []JsonItem {
	items := []JsonItem{}

	for _, file := range fg.Options.Files {
		items = append(items, JsonItem{
			Id:          convert.FormatLink(fg.Options.SiteUrl, file.Data.Slug),
			Title:       file.Data.Title,
			Url:         convert.FormatLink(fg.Options.SiteUrl, file.Data.Slug),
			ContentHtml: file.Html,
			ContentText: file.Text,
			Summary:     file.Data.Description,
			Published:   jsonDate(file.Data.CreatedDate),
			Updated:     jsonDate(file.Data.UpdatedDate),
			Authors: []JsonAuthor{{
				Name: convert.DefaultTo(file.Data.Author, fg.Options.Author).(string),
				Url:  fg.Options.SiteUrl,
			}},
		})
	}

	return items
}

func (fg *FeedGenerator) Generate() {
	atom, err := xml.Marshal(&fg.Atom)

	if err != nil {
		panic(err)
	}

	rss, err := xml.Marshal(&fg.Rss)

	if err != nil {
		panic(err)
	}
	r := regexp.MustCompile("></atom:link>")
	rssStr := r.ReplaceAllString(string(rss), "/>")

	json, err := json.Marshal(&fg.Json)

	if err != nil {
		panic(err)
	}

	file.WriteFile(path.Join(fg.Options.RootPath, fg.Options.RelativeFeedPath, "feed.atom"), xml.Header+string(atom))
	file.WriteFile(path.Join(fg.Options.RootPath, fg.Options.RelativeFeedPath, "feed.rss"), xml.Header+rssStr)
	file.WriteFile(path.Join(fg.Options.RootPath, fg.Options.RelativeFeedPath, "feed.json"), string(json))
}

func NewFeedGenerator(options *FeedGeneratorOptions) *FeedGenerator {
	now := time.Now()
	files := options.Files
	lastPublishedDate := findLastPublished(files)

	fg := &FeedGenerator{
		Options: options,
	}

	feedLink := convert.FormatLink(options.SiteUrl, options.RelativeFeedPath)

	atom := &AtomFeed{
		NS:    "http://www.w3.org/2005/Atom",
		Title: options.Title,
		Link: []AtomLink{{
			Rel:  "self",
			Href: feedLink + "feed.atom",
		}, {
			Href: feedLink,
		}},
		Id:      feedLink,
		Updated: atomDate(lastPublishedDate),
		Icon:    options.IconSm,
		Logo:    options.Logo2x1,
		Author: &AtomAuthor{
			Name: options.Author,
		},
		Contributor: &AtomAuthor{
			Name: options.Author,
		},
		Entries: fg.getAtomItems(),
	}

	rss := &RssFeed{
		Version: "2.0",
		Atom:    "http://www.w3.org/2005/Atom",
		Dc:      "http://purl.org/dc/elements/1.1/",
		Channel: RssChannel{
			Title:         options.Title,
			Link:          feedLink,
			Description:   options.Description,
			Language:      convert.DefaultTo(options.Language, "en-US").(string),
			Copyright:     options.Copyright,
			PubDate:       rssDate(lastPublishedDate),
			LastBuildDate: rssDate(now),
			Docs:          "https://cyber.harvard.edu/rss/rss.html",
			Items:         fg.getRssItems(),
			Image: &RssImage{
				URL:   options.IconLg,
				Title: options.Title,
				Link:  options.SiteUrl,
			},
			AtomLink: RssAtomLink{
				Href: feedLink + "feed.rss",
				Rel:  "self",
				Type: "application/rss+xml",
			},
		},
	}

	json := &JsonFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       options.Title,
		HomePageUrl: options.SiteUrl,
		FeedUrl:     feedLink + "feed.json",
		Description: options.Description,
		Authors: []JsonAuthor{{
			Name: options.Author,
			Url:  options.SiteUrl,
		}},
		Icon:     options.IconLg,
		Favicon:  options.IconMd,
		Language: convert.DefaultTo(options.Language, "en-US").(string),
		Items:    fg.getJsonItems(),
	}

	fg.Atom = atom
	fg.Rss = rss
	fg.Json = json
	return fg
}
