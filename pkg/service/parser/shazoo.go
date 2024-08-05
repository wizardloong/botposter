package parser

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Shazoo struct {
}

func (parser *Shazoo) FetchArticleData(url string) (string, string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	content := doc.Find("section.Entry__content").Text()
	title := doc.Find("article section h1").Text()
	imageURL, exists := doc.Find("figure img.w-full").Attr("src")
	if !exists {
		return "", "", "", fmt.Errorf("image not found")
	}

	return content, title, imageURL, nil
}

func (parser *Shazoo) DownloadImage(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
