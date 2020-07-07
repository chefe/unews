package srf

import (
	"github.com/antchfx/htmlquery"
	"github.com/chefe/unews/helper"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"log"
	"strings"
	"sync"
	"time"
)

func GetFeed() (*feeds.Feed, error) {
	fp := gofeed.NewParser()
	source, err := fp.ParseURL("https://www.srf.ch/news/bnf/rss/1646")
	if err != nil {
		log.Fatal(err)
	}

	feed := &feeds.Feed{
		Title:       source.Title,
		Link:        &feeds.Link{Href: source.Link},
		Description: source.Description,
		Created:     time.Now(),
	}

	ch := make(chan *feeds.Item)
	go func(source *gofeed.Feed, ch chan *feeds.Item) {
		var wg sync.WaitGroup

		for _, item := range source.Items {
			wg.Add(1)
			go createItem(item, ch, &wg)
		}

		wg.Wait()
		close(ch)
	}(source, ch)

	for item := range ch {
		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}

func createItem(source *gofeed.Item, ch chan *feeds.Item, wg *sync.WaitGroup) {
	doc, err := htmlquery.LoadURL(source.Link)
	if err != nil {
		log.Fatal(err)
		return
	}

	id := getContentAttrOf(doc, "//meta[@itemprop='identifier']")
	lead := getContentAttrOf(doc, "//meta[@name='description']")
	title := getInnerTextOf(doc, "//span[@class='article-title__text']")
	titleHeader := getInnerTextOf(doc, "//span[@class='article-title__overline']")

	createdAsStr := getContentAttrOf(doc, "//meta[@itemprop='datePublished']")
	created, _ := time.Parse(time.RFC3339, createdAsStr)

	updatedAsStr := getContentAttrOf(doc, "//meta[@itemprop='dateModified']")
	updated, _ := time.Parse(time.RFC3339, updatedAsStr)

	item := &feeds.Item{
		Id:          id,
		Title:       titleHeader + " — " + title,
		Link:        &feeds.Link{Href: source.Link},
		Description: strings.TrimSpace(lead),
		Content:     getContent(doc, source.Link),
		Created:     created,
		Updated:     updated,
	}

	ch <- item
	wg.Done()
}

func getContentAttrOf(doc *html.Node, xpath string) string {
	node := htmlquery.FindOne(doc, xpath)
	value, _ := getAttr(node.Attr, "content")
	return value
}

func getInnerTextOf(doc *html.Node, xpath string) string {
	node := htmlquery.FindOne(doc, xpath)
	return htmlquery.InnerText(node)
}

func getAttr(attr []html.Attribute, name string) (string, bool) {
	for _, a := range attr {
		if a.Key == name {
			return a.Val, true
		}
	}

	return "", false
}

func getClassOf(node *html.Node) string {
	class, _ := getAttr(node.Attr, "class")
	return strings.TrimSpace(class)
}

func getContent(doc *html.Node, url string) string {
	contentNode := htmlquery.FindOne(doc, "//div[@class='article-content']")
	content := ""

	for n := contentNode.FirstChild; n != nil; n = n.NextSibling {
		switch strings.TrimSpace(n.Data) {
		case "p":
			content += htmlquery.OutputHTML(n, true)
		case "ul":
			content += "<p>" + htmlquery.OutputHTML(n, true) + "</p>"
		case "h2", "h3":
			content += "<h3>" + htmlquery.OutputHTML(n, false) + "</h3>"
		case "div":
			content += handleDiv(n, url)
		case "a":
			// ignore
		case "googleoff: all":
			// ignore
		case "googleon: all":
			// ignore
		case "":
			// ignore
		default:
			log.Print("Unsupported element: " + n.Data + "\n -> " + url)
		}
	}

	return content
}

func handleDiv(node *html.Node, url string) string {
	html := ""

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		switch strings.TrimSpace(n.Data) {
		case "figure":
			html += createFigure(n)
		case "blockquote":
			html += createBlockquote(n)
		case "div":
			class := getClassOf(n)

			switch {
			case strings.HasPrefix(class, "expandable-box"):
				// ignore
			case strings.HasPrefix(class, "embed-inline"):
				// ignore
			case strings.HasPrefix(class, "js-iframe-modal-caller"):
				// ignore
			case strings.HasPrefix(class, "carousel-container"):
				// ignore
			case strings.HasPrefix(class, "related-items-list"):
				// ignore
			case strings.HasPrefix(class, "media-window-container"):
				// ignore
			case strings.HasPrefix(class, "linkbox"):
				// ignore
			case class == "":
				// ignore
			default:
				log.Print("Unsupported class for div-tag: " + class + "\n -> " + url)
			}
		case "a":
			class := getClassOf(n)

			switch class {
			case "js-media":
				// ignore
			default:
				log.Print("Unsupported class for anchor-tag: " + class + "\n -> " + url)
			}
		case "googleoff: all":
			// ignore
		case "googleon: all":
			// ignore
		case "":
			// ignore
		default:
			log.Print("Unsupported div element: " + n.Data + "\n -> " + url)
		}
	}

	return html
}

func createFigure(node *html.Node) string {
	srcNode := htmlquery.FindOne(node, "//img/@src")
	src, _ := helper.GetImageAsBase64URL(htmlquery.InnerText(srcNode))
	html := "<img src=\"" + src + "\">"

	captionNode := htmlquery.FindOne(node, "//span[@class='media-caption__description']")
	if captionNode != nil {
		caption := strings.TrimSpace(htmlquery.InnerText(captionNode))
		html += "<figcaption>" + caption + "</figcaption>"
	}

	return "<figure>" + html + "</figure>"
}

func createBlockquote(node *html.Node) string {
	quoteNode := htmlquery.FindOne(node, "//span[@class='blockquote__text']")
	quote := strings.TrimSpace(htmlquery.InnerText(quoteNode))
	html := "<p>" + quote + "</p>"

	authorNode := htmlquery.FindOne(node, "//span[@class='blockquote__author']")
	if authorNode != nil {
		author := strings.TrimSpace(htmlquery.InnerText(authorNode))
		html += "<footer>" + author + "</footer>"
	}

	return "<blockquote>" + html + "</blockquote>"
}
