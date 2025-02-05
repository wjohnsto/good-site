package rss

import (
	"time"

	"good.site/pkg/services"
)

func atomDate(dt time.Time) string {
	return dt.Format(time.RFC3339)
}

func jsonDate(dt time.Time) string {
	return dt.Format(time.RFC3339)
}

func rssDate(dt time.Time) string {
	return dt.Format("02 Jan 2006 15:04 -0700")
}

func findLastPublished(files []*services.MarkdownFile) time.Time {
	if len(files) == 0 {
		return time.Now()
	}

	lastPublished := files[0].Data.CreatedDate

	for _, file := range files {
		if file.Data.CreatedDate.After(lastPublished) {
			lastPublished = file.Data.CreatedDate
		}
	}

	return lastPublished
}
