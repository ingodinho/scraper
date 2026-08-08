package internal

import (
	"testing"
)

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
		actual := getHeadingFromHtml(tc.inputHtml)
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
			actual := getFirstParagraphFromHtml(tc.inputHtml)
			if actual != tc.expected {
				t.Errorf("Test: %v - Actual: %v not as expected: %v", i, actual, tc.expected)
			}
		})
	}
}
