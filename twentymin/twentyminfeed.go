package twentymin

import (
	"encoding/json"
	"github.com/chefe/unews/helper"
	"github.com/chefe/unews/twentymin/crawler"
	"github.com/chefe/unews/twentymin/json/index"
	"github.com/chefe/unews/twentymin/json/post"
	"github.com/gorilla/feeds"
	"log"
	"strings"
	"sync"
	"time"
)

func GetFeed() (*feeds.Feed, error) {
	inner, err := crawler.GetNextDataJSONFromPage("https://www.20min.ch")
	if err != nil {
		return nil, err
	}

	var n index.IndexPageJSON
	err = json.Unmarshal([]byte(inner), &n)

	feed := &feeds.Feed{
		Title:       n.Props.Data.Content.MetaTitle,
		Link:        &feeds.Link{Href: "https://www.20min.ch"},
		Description: n.Props.Data.Content.MetaDescription,
		Created:     time.Now(),
	}

	ch := make(chan *feeds.Item)
	go func(n index.IndexPageJSON, ch chan *feeds.Item) {
		var wg sync.WaitGroup

		for _, c := range n.Props.Data.Content.Elements {
			for _, p := range c.Elements {
				if p.Type == "articles" {
					wg.Add(1)
					go createItem(p, ch, &wg)
				}
			}
		}

		wg.Wait()
		close(ch)
	}(n, ch)

	for item := range ch {
		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}

func createItem(p index.PostElementJSON, ch chan *feeds.Item, wg *sync.WaitGroup) {
	updated, _ := time.Parse(time.RFC3339, p.Content.Updated)
	created, _ := time.Parse(time.RFC3339, p.Content.Published)
	url := "https://www.20min.ch" + p.Content.URL

	item := &feeds.Item{
		Id:          p.ID,
		Title:       p.Content.TitleHeader + " — " + p.Content.Title,
		Link:        &feeds.Link{Href: url},
		Description: strings.TrimSpace(p.Content.Lead),
		Content:     getContent(p.Content.URL),
		Created:     created,
		Updated:     updated,
	}

	ch <- item
	wg.Done()
}

func getContent(url string) string {
	url = "https://www.20min.ch" + url

	inner, err := crawler.GetNextDataJSONFromPage(url)
	if err != nil {
		log.Fatal("Error while processing " + url)
		log.Fatal(err)
		return "<p>Inhalt konnte nicht geladen werden!</p>"
	}

	var n post.PostPageJSON
	err = json.Unmarshal([]byte(inner), &n)
	if err != nil {
		log.Fatal("Error while processing " + url)
		log.Fatal(err)
		return "<p>Inhalt konnte nicht geladen werden!</p>"
	}

	var content string

	for _, p := range n.Props.Data.Content.Article.Elements {
		switch p.Type {
		case "title-header":
			// ignore
		case "title":
			// ignore
		case "lead":
			content += "<p><strong><em>" + p.HTMLText + "</em></strong></p>"
		case "textBlockArray":
			content += "<p>"

			for _, e := range p.Items {
				switch e.Type {
				case "htmlTextItem":
					content += helper.RemoveLinks(e.HTMLText)
				case "internalLink":
					content += e.HTMLText
				default:
					log.Fatal("Error while processing " + url)
					log.Fatal("Unsupported block element typ: " + e.Type)
				}
			}

			content += "</p>"
		case "crosshead":
			content += "<h3>" + p.HTMLText + "</h3>"
		case "unordered-list":
			content += "<p>" + p.HTMLText + "</p>"
		case "quote":
			content += "<blockquote><p>" + p.Quote + "</p><footer>"
			content += p.Author + "</footer></blockquote>"
		case "container":
			// ignore
		case "slideshow":
			// ignore
		case "agencies":
			// ignore
		case "image":
			content += createFigureHTML(p.Image)
		case "videocms":
			content += createVideoHTML(p)
		case "authors":
			// ignore
		case "ad":
			// ignore
		case "embed":
			// ignore
		case "publishDate":
			// ignore
		case "footer":
			// ignore
		default:
			log.Fatal("Error while processing " + url)
			log.Fatal("Unsupported article element typ: " + p.Type)
		}
	}

	content += "<p><small>Quelle: <a href=\"" + url + "\">20min.ch</a></small></p>"
	return content
}

func createFigureHTML(image post.ImageJSON) string {
	html := ""

	url := image.Variants.Small.Src
	src, err := helper.GetImageAsBase64URL(url)
	if err == nil {
		html += "<figure><img src=\"" + src + "\"><figcaption>"
		html += image.Caption.Text + "</figcaption></figure>"
	}

	return html
}

func createVideoHTML(p post.ArticleElementJSON) string {
	src := p.Content.Elements[0].VideoURL
	poster, _ := helper.GetImageAsBase64URL(p.Content.Elements[0].VideoThumbnail)

	html := "<video controls preload=\"metadata\" poster=\"" + poster + "\">"
	html += "<source src=\"" + src + "\" type=\"video/mp4\">"
	html += "Dein Browser unterstützt keine Videos.</video>"

	return html
}
