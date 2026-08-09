package internal

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func ExtractPageData(html, pageUrl string) PageData {
	url, err := url.Parse(pageUrl)
	if err != nil {
		return PageData{}
	}

	heading := getHeadingFromHTML(html)

	paragraph := getFirstParagraphFromHTML(html)

	outgoingLinks, err := getURLsFromHTML(html, url)
	if err != nil {
		fmt.Printf("error getting urls from html: %v", err)
	}

	imageURLs, err := getImagesFromHTML(html, url)
	if err != nil {
		fmt.Printf("error getting images from html: %v", err)
	}

	return PageData{
		URL:            url.String(),
		Heading:        heading,
		FirstParagraph: paragraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	heading := strings.TrimSpace(doc.Find("h1").First().Text())
	if heading == "" {
		heading = strings.TrimSpace(doc.Find("h2").First().Text())
	}

	return heading
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	p := strings.TrimSpace(doc.Find("main").First().Find("p").First().Text())
	if p == "" {
		p = strings.TrimSpace(doc.Find("p").First().Text())
	}

	return p
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	result := []string{}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		if strings.HasPrefix(href, "/") {
			href = baseURL.String() + href
		}

		result = append(result, href)
	})

	return result, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return []string{}, err
	}

	result := []string{}

	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, _ := s.Attr("src")

		if strings.HasPrefix(src, "/") {
			src = baseURL.String() + src
		}

		result = append(result, src)
	})

	return result, nil
}
