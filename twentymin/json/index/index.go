package index

type IndexPageJSON struct {
	Props PropsJSON
}

type PropsJSON struct {
	Data DataJSON
}

type DataJSON struct {
	Content ContentJSON
}

type ContentJSON struct {
	Elements        []CategoryElementJSON
	MetaTitle       string
	MetaDescription string
}

type CategoryElementJSON struct {
	Elements []PostElementJSON
	Title    string
}

type PostElementJSON struct {
	Content PostContentJSON
	Type    string
	ID      string
}

type PostContentJSON struct {
	Lead        string
	Title       string
	TitleHeader string
	Published   string
	Updated     string
	URL         string
}
