package internal

import (
	"net/url"
	"strings"
)

func normalizeUrl(inputUrl string) (string, error) {
	parsed, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}

	path, _ := strings.CutSuffix(parsed.Path, "/")

	outputUrl, err := url.JoinPath(parsed.Host, path)
	if err != nil {
		return "", err
	}

	return outputUrl, nil
}
