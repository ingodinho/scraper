package internal

import "testing"

func TestNormalizeUrl(t *testing.T) {
	tests := []struct {
		name     string
		inputUrl string
		expected string
	}{{
		name:     "scheme_with_trailing_slash",
		inputUrl: "https://www.boot.dev/blog/path/",
		expected: "www.boot.dev/blog/path",
	},
		{
			name:     "scheme_without_trailing_slash",
			inputUrl: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "http_with_trailing_slash",
			inputUrl: "http://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "http_without_trailing_slash",
			inputUrl: "http://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeUrl(tc.inputUrl)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
