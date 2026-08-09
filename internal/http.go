package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetHTML(rawURL string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v\n", err)
	}

	req.Header.Add("User-Agent", "BootCrawler/1.0")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %v\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return "", fmt.Errorf("bad status code: %v\n", res.StatusCode)
	}

	contentType := res.Header.Get("content-type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("response does not have content type text/html: %v \n", contentType)
	}

	html, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	
	return string(html), nil
}
