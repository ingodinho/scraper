package internal

import (
	"net/url"
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  PageData
	}{{
		name:     "one of each",
		inputURL: "https://crawler-test.com",
		inputBody: `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`,
		expected: PageData{
			URL:            "https://crawler-test.com",
			Heading:        "Test Title",
			FirstParagraph: "This is the first paragraph.",
			OutgoingLinks:  []string{"https://crawler-test.com/link1"},
			ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
		},
	},
		{
			name:      "empty page",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body></body></html>`,
			expected: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "",
				FirstParagraph: "",
				OutgoingLinks:  []string{},
				ImageURLs:      []string{},
			},
		},
		{
			name:     "multiple links and images",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
        <h1>Multi Test</h1>
        <p>First paragraph here.</p>
        <p>Second paragraph should be ignored.</p>
        <a href="/link1">Link 1</a>
        <a href="https://crawler-test.com/link2">Link 2</a>
        <a href="https://other.com/link3">Link 3</a>
        <img src="/image1.jpg" alt="Image 1">
        <img src="https://other.com/image2.png">
    </body></html>`,
			expected: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Multi Test",
				FirstParagraph: "First paragraph here.",
				OutgoingLinks: []string{
					"https://crawler-test.com/link1",
					"https://crawler-test.com/link2",
					"https://other.com/link3",
				},
				ImageURLs: []string{
					"https://crawler-test.com/image1.jpg",
					"https://other.com/image2.png",
				},
			},
		},
		{
			name:     "main tag takes priority",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
        <h1>Outer Heading</h1>
        <p>Outer paragraph.</p>
        <main>
            <p>Main paragraph.</p>
        </main>
    </body></html>`,
			expected: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Outer Heading",
				FirstParagraph: "Main paragraph.",
				OutgoingLinks:  []string{},
				ImageURLs:      []string{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractPageData(tc.inputBody, tc.inputURL)
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Test: %v - Expected: %v - Actual: %v", tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetHeadingFromHtml(t *testing.T) {
	tests := []struct {
		name      string
		inputHtml string
		expected  string
	}{
		{
			name: "get h1 tag",
			inputHtml: `<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
    <main>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
  </body>
			</html>`,
			expected: "Welcome to Boot.dev",
		},
		{
			name: "get h2 tag",
			inputHtml: `<html>
  <body>
    <h2>Welcome to Boot.dev h2</h2>
    <main>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
  </body>
			</html>`,
			expected: "Welcome to Boot.dev h2",
		},
		{
			name: "get default string",
			inputHtml: `<html>
  <body>
    <main>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
  </body>
			</html>`,
			expected: "",
		},
	}

	for i, tc := range tests {
		actual := getHeadingFromHTML(tc.inputHtml)
		if actual != tc.expected {
			t.Errorf("Test: %v - Actual: %v does not match expected: %v", i, actual, tc.expected)
		}
	}
}

func TestGetParagraphFromHtml(t *testing.T) {
	tests := []struct {
		name      string
		inputHtml string
		expected  string
	}{{
		name: "p inside main",
		inputHtml: `<html>
  <body>
    <main>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
  </body>
	</html>`,
		expected: "Learn to code by building real projects.",
	},
		{
			name: "p outside main",
			inputHtml: `<html>
  <body>
    <p>Learn to code by building real projects.</p>
    <p>This is the second paragraph.</p>
    <main>
      <div>No paragraphs here.</div>
    </main>
  </body>
		</html>`,
			expected: "Learn to code by building real projects.",
		},
		{
			name: "p inside and outside main",
			inputHtml: `<html>
  <body>
    <p>Paragraph outside main.</p>
    <main>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
    <p>Another paragraph outside main.</p>
  </body>
		</html>`,
			expected: "Learn to code by building real projects.",
		}, {
			name: "no p at all",
			inputHtml: `<html>
  <body>
    <main>
      <div>Learn to code by building real projects.</div>
      <span>No paragraph elements in this document.</span>
    </main>
  </body>
		</html>`,
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.inputHtml)
			if actual != tc.expected {
				t.Errorf("Test: %v - Actual: %v not as expected: %v", i, actual, tc.expected)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputUrl  string
		inputBody string
		expected  []string
	}{
		{
			name:      "one url",
			inputUrl:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com"},
		},
		{
			name:      "no anchor tags",
			inputUrl:  "https://crawler-test.com",
			inputBody: `<html><body><p>No links here</p><div><span>Boot.dev</span></div></body></html>`,
			expected:  []string{},
		},
		{
			name:     "three urls",
			inputUrl: "https://crawler-test.com",
			inputBody: `<html><body>
				<a href="https://crawler-test.com/one"><span>One</span></a>
				<a href="/two">Two</a>
				<a href="https://other.com/three">Three</a>
			</body></html>`,
			expected: []string{
				"https://crawler-test.com/one",
				"https://crawler-test.com/two",
				"https://other.com/three",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputUrl)
			if err != nil {
				t.Errorf("error parsing url: %v", err)
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Errorf("error getting url from html: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("test: %v - expected: %v - actual: %v", i, tc.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputUrl  string
		inputBody string
		expected  []string
	}{{
		name:      "one image",
		inputUrl:  "https://crawler-test.com",
		inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
		expected:  []string{"https://crawler-test.com/logo.png"},
	},
		{
			name:      "no image tags",
			inputUrl:  "https://crawler-test.com",
			inputBody: `<html><body><p>No images here</p><a href="/page">Link</a></body></html>`,
			expected:  []string{},
		},
		{
			name:     "three images",
			inputUrl: "https://crawler-test.com",
			inputBody: `<html><body>
			<img src="/logo.png" alt="Logo">
			<img src="https://crawler-test.com/banner.jpg" alt="Banner">
			<div><img src="https://other.com/icon.svg"></div>
		</body></html>`,
			expected: []string{
				"https://crawler-test.com/logo.png",
				"https://crawler-test.com/banner.jpg",
				"https://other.com/icon.svg",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputUrl)
			if err != nil {
				t.Errorf("error parsing url: %v", err)
			}

			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Errorf("error getting url from html: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("test: %v - expected: %v - actual: %v", i, tc.expected, actual)
			}
		})
	}
}
