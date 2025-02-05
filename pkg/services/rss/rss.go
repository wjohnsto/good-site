package rss

// RssFeed represents the RSS feed channel (RSS 2.0 specification)
type RssFeed struct {
	XMLName struct{}   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Atom    string     `xml:"xmlns:atom,attr"`
	Dc      string     `xml:"xmlns:dc,attr"`
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language,omitempty"`
	Copyright   string `xml:"copyright,omitempty"`
	// RFC822, Latest post created date
	PubDate string `xml:"pubDate,omitempty"`
	// RFC822, Last site build
	LastBuildDate string `xml:"lastBuildDate,omitempty"`
	Generator     string `xml:"generator,omitempty"`
	Docs          string `xml:"docs,omitempty"`
	// <email> (name)
	ManagingEditor string `xml:"managingEditor,omitempty"`
	// <email> (name)
	WebMaster string      `xml:"webMaster,omitempty"`
	Image     *RssImage   `xml:"image,omitempty"`
	AtomLink  RssAtomLink `xml:"atom:link"`
	Items     []RssItem   `xml:"item"`
}

type RssAtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

// RssImage represents the image associated with the RSS channel
type RssImage struct {
	URL    string `xml:"url"`
	Title  string `xml:"title"`
	Link   string `xml:"link"`
	Width  int    `xml:"width,omitempty"`
	Height int    `xml:"height,omitempty"`
}

// RssItem represents an individual item (entry) in the RSS feed
type RssItem struct {
	Title       string             `xml:"title"`
	Link        string             `xml:"link"`
	Description RssItemDescription `xml:"description"`
	// <email> (name)
	Author   string `xml:"dc:creator,omitempty"`
	Category string `xml:"category,omitempty"`
	// RFC822
	PubDate   string        `xml:"pubDate,omitempty"`
	Guid      string        `xml:"guid,omitempty"`
	Enclosure *RssEnclosure `xml:"enclosure,omitempty"`
	Source    *RssSource    `xml:"source,omitempty"`
}

type RssItemDescription struct {
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr,omitempty"`
}

// RssEnclosure represents the enclosure of an item (e.g., media file)
type RssEnclosure struct {
	URL    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

// RssSource represents the source of an item (optional field in RSS)
type RssSource struct {
	Url string `xml:"url,attr"`
}
