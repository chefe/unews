package crawler

import (
	"bytes"
	"encoding/base64"
	"errors"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func hasID(attr []html.Attribute, id string) bool {
	for _, a := range attr {
		if a.Key == "id" && a.Val == id {
			return true
		}
	}

	return false
}

func findNextDataNode(node *html.Node) (*html.Node, error) {
	if node.Type == html.ElementNode && node.Data == "script" && hasID(node.Attr, "__NEXT_DATA__") {
		return node, nil
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		n, _ := findNextDataNode(child)
		if n != nil {
			return n, nil
		}
	}

	return nil, errors.New("Node not found")
}

func getNodeAsHTML(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func getJSONFromNextData(body string) string {
	trimStart := "<script id=\"__NEXT_DATA__\" type=\"application/json\">"
	trimEnd := "</script>"
	inner := body[len(trimStart) : len(body)-len(trimEnd)]
	return inner
}

func GetNextDataJSONFromPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Failed to load the url")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Failed to read the body")
	}

	doc, err := html.Parse(strings.NewReader(string(bytes)))
	if err != nil {
		return "", errors.New("Failed to parse the body")
	}

	dataNode, err := findNextDataNode(doc)
	if err != nil {
		return "", err
	}

	htm := getNodeAsHTML(dataNode)
	return getJSONFromNextData(htm), nil
}

func GetImageAsBase64URL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Failed to load the url")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	mediaType := resp.Header["Content-Type"][0]
	image := base64.StdEncoding.EncodeToString(bytes)
	return "data:" + mediaType + ";base64, " + image, nil
}
