package internal

import (
	"fmt"
	"net/url"
	"strings"
)

func CrawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) error {
	fmt.Printf("crawling %v....\n", rawCurrentURL)
	if !strings.HasPrefix(rawCurrentURL, rawBaseURL) {
		return fmt.Errorf("rawCurrentUrl: %v has not the same prefix as rawBaseUrl: %v\n", rawCurrentURL, rawBaseURL)
	}

	normalizedCurrent, err := normalizeUrl(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing the current url %v: %v \n", normalizedCurrent, err)
		return err
	}

	_, exists := pages[normalizedCurrent]
	if exists {
		pages[normalizedCurrent] += 1
		return nil
	} else {
		pages[normalizedCurrent] = 1
	}

	html, err := GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting html: %v \n", err)
		return err
	}

	// fmt.Printf("html from %v: %v \n", rawCurrentURL, html)

	parsedUrl, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing base url: %v\n", err)
		return err
	}

	urls, err := getURLsFromHTML(html, parsedUrl)
	if err != nil {
		fmt.Printf("error getting urls from html: %v\n", err)
		return err
	}

	for _, url := range urls {
		err := CrawlPage(rawBaseURL, url, pages)
		if err != nil {
			fmt.Printf("error crawlpage %v \n", err)
		}
	}

	return nil
}
