package rss

// AtomFeed represents the Atom feed (Atom 1.0 specification)
type AtomFeed struct {
	XMLName struct{}   `xml:"feed"`
	NS      string     `xml:"xmlns,attr"`
	Title   string     `xml:"title"`
	Link    []AtomLink `xml:"link"`
	Id      string     `xml:"id"`
	Logo    string     `xml:"logo"`
	Icon    string     `xml:"icon"`
	// RFC3339
	Updated     string      `xml:"updated"`
	Author      *AtomAuthor `xml:"author,omitempty"`
	Contributor *AtomAuthor `xml:"contributor,omitempty"`
	Entries     []AtomEntry `xml:"entry"`
}

// AtomLink represents the link element in Atom, which can have attributes like href
type AtomLink struct {
	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr,omitempty"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
}

// AtomAuthor represents the author element in Atom, typically includes name and email
type AtomAuthor struct {
	Name  string `xml:"name"`
	Email string `xml:"email,omitempty"`
}

// AtomEntry represents an individual entry (item) in the Atom feed
type AtomEntry struct {
	Title string     `xml:"title"`
	Link  []AtomLink `xml:"link"`
	Id    string     `xml:"id"`
	// RFC3339
	Updated    string         `xml:"updated"`
	Summary    string         `xml:"summary,omitempty"`
	Content    AtomContent    `xml:"content,omitempty"`
	Author     *AtomAuthor    `xml:"author,omitempty"`
	Categories []AtomCategory `xml:"category,omitempty"`
	// RFC3339
	Published string      `xml:"published,omitempty"`
	Source    *AtomSource `xml:"source,omitempty"`
}

type AtomContent struct {
	Html string `xml:",chardata"`
	Type string `xml:"type,attr"`
}

// AtomCategory represents a category associated with an Atom entry
type AtomCategory struct {
	Term   string `xml:"term,attr"`
	Scheme string `xml:"scheme,attr,omitempty"`
	Label  string `xml:"label,attr,omitempty"`
}

// AtomSource represents the source feed of an Atom entry
type AtomSource struct {
	Id    string `xml:"id"`
	Title string `xml:"title"`
	// RFC3339
	Updated string     `xml:"updated"`
	Link    []AtomLink `xml:"link"`
}
