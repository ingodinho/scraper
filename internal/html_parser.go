package internal

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHtml(html string) string {
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


func getFirstParagraphFromHtml(html string) string {
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
