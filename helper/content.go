package helper

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

func RemoveLinks(html string) string {
	re := regexp.MustCompile(`<a href[^>]*>(.*)<\/a>`)
	return string(re.ReplaceAll([]byte(html), []byte("$1")))
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
