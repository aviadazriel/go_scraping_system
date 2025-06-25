package models

// ParserConfig represents the configuration for parsing HTML
// This is a simplified version for the API Gateway
type ParserConfig struct {
	TitleSelector   string            `json:"title_selector,omitempty"`
	ContentSelector string            `json:"content_selector,omitempty"`
	AuthorSelector  string            `json:"author_selector,omitempty"`
	DateSelector    string            `json:"date_selector,omitempty"`
	ImageSelector   string            `json:"image_selector,omitempty"`
	PriceSelector   string            `json:"price_selector,omitempty"`
	CustomSelectors map[string]string `json:"custom_selectors,omitempty"`
	ExtractMetadata bool              `json:"extract_metadata,omitempty"`
	ExtractLinks    bool              `json:"extract_links,omitempty"`
	ExtractImages   bool              `json:"extract_images,omitempty"`
	RemoveScripts   bool              `json:"remove_scripts,omitempty"`
	RemoveStyles    bool              `json:"remove_styles,omitempty"`
	CleanHTML       bool              `json:"clean_html,omitempty"`
}
