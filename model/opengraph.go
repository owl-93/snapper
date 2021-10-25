package model

type SnapperResult struct {
	Url string`json:"url"`
	Title string`json:"title"`
	Description string`json:"description"`
	Image string`json:"image"`
	Type string`json:"type"`
	Locale string`json:"locale"`
}

const (
	UrlKey = "og:url"
	TitleKey = "og:title"
	DescKey = "og:description"
	ImageKey = "og:image"
	TypeKey = "og:type"
	LocalKey = "og:locale"
)
