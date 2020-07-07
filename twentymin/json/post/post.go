package post

type PostPageJSON struct {
	Props PropsJSON `json:"props"`
}

type PropsJSON struct {
	Data DataJSON `json:"data"`
}

type DataJSON struct {
	Content ContentJSON `json:"content"`
}

type ContentJSON struct {
	Article ArticleJSON `json:"article"`
}

type ArticleJSON struct {
	Elements []ArticleElementJSON `json:"elements"`
}

type ArticleElementJSON struct {
	Type     string     `json:"type"`
	HTMLText string     `json:"htmlText"`
	Items    []ItemJSON `json:"items"`

	// for quote type
	Quote  string `json:"quote"`
	Author string `json:"author"`

	// for internalLink type
	URL string `json:"url"`

	// for slideshow type
	Slideshow SlideshowJSON `json:"slideshow"`

	// for image type
	Image ImageJSON `json:"image"`

	// for videocms type
	VideoThumbnail string      `json:"thumbnail"`
	VideoURL       string      `json:"url_high"`
	Content        ArticleJSON `json:"content"`

	// for container type
	Elements []ArticleElementJSON `json:"elements"`
}

type ImageJSON struct {
	Caption  ImageCaptionJSON  `json:"caption"`
	Variants ImageVariantsJSON `json:"variants"`
}

type ImageCaptionJSON struct {
	Text string `json:"text"`
}

type ImageVariantsJSON struct {
	Big   ImageVariantJSON `json:"big"`
	Small ImageVariantJSON `json:"small"`
}

type ImageVariantJSON struct {
	Src string `json:"src"`
}

type SlideshowJSON struct {
	Slides []SlideJSON `json:"slides"`
}

type SlideJSON struct {
	Image ImageJSON `json:"image"`
	Type  string    `json:"type"`
}

type ItemJSON struct {
	Type     string
	HTMLText string
}
