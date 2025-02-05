package rss

// JsonFeed represents the JSON Feed (version 1.1 specification)
type JsonFeed struct {
	Version     string                 `json:"version"`       // version of the JSON Feed format
	Title       string                 `json:"title"`         // Title of the feed
	HomePageUrl string                 `json:"home_page_url"` // URL to the homepage
	FeedUrl     string                 `json:"feed_url"`      // URL to the feed itself
	Description string                 `json:"description,omitempty"`
	Icon        string                 `json:"icon,omitempty"`       // URL to the feed's icon
	Favicon     string                 `json:"favicon,omitempty"`    // URL to the favicon of the feed
	Authors     []JsonAuthor           `json:"authors,omitempty"`    // Author of the feed
	Items       []JsonItem             `json:"items"`                // Items contained in the feed
	Language    string                 `json:"language,omitempty"`   // Language of the feed
	Extensions  map[string]interface{} `json:"extensions,omitempty"` // Custom extensions to the feed
}

// JsonAuthor represents the author of the JSON Feed
type JsonAuthor struct {
	Name string `json:"name"`          // Name of the author
	Url  string `json:"url,omitempty"` // URL to the author's homepage (optional)
}

// JsonItem represents an individual item (entry) in the JSON Feed
type JsonItem struct {
	Id          string `json:"id"`                     // Unique ID of the item
	Title       string `json:"title"`                  // Title of the item
	Url         string `json:"url"`                    // URL to the item
	ContentHtml string `json:"content_html,omitempty"` // HTML content for the item
	ContentText string `json:"content_text,omitempty"` // Plain text content for the item
	Summary     string `json:"summary,omitempty"`      // Summary of the item
	Image       string `json:"image,omitempty"`        // Summary of the item
	// RFC3339
	Published string `json:"date_published,omitempty"` // Published timestamp for the item
	// RFC3339
	Updated     string                 `json:"date_modified,omitempty"` // Last updated timestamp for the item
	Authors     []JsonAuthor           `json:"authors,omitempty"`       // Author of the item
	Tags        []string               `json:"tags,omitempty"`          // Tags associated with the item
	Attachments []JsonAttachment       `json:"attachments,omitempty"`   // Attachments for the item
	Language    string                 `json:"language,omitempty"`      // Language of the item
	Extensions  map[string]interface{} `json:"extensions,omitempty"`    // Custom extensions for the item
}

// JsonAttachment represents an attachment (e.g., media file) in a JSON Feed item
type JsonAttachment struct {
	Url      string `json:"url"`             // URL of the attachment
	MimeType string `json:"mime_type"`       // MIME type of the attachment (e.g., "image/jpeg")
	Size     int64  `json:"size,omitempty"`  // Size of the attachment in bytes
	Title    string `json:"title,omitempty"` // Title of the attachment
}
