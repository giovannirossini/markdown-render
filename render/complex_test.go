package render

import (
	"testing"
)

func TestRender_ComplexDocument(t *testing.T) {
	markdown := `# Main Title

This is a paragraph with **bold** and *italic* text.

## Section

- Item 1
- Item 2
  - Nested item

` + "`inline code`" + ` and [a link](http://example.com)

` + "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```" + `

| Col1 | Col2 |
|------|------|
| A    | B    |
`

	result := RenderToString(markdown)

	want := []string{
		"# ", "Main Title",
		"bold", "italic",
		"## ", "Section",
		"Item 1", "Item 2", "Nested item",
		"inline code",
		"a link", "http://example.com",
		"func main()", "fmt.Println",
		"Col1", "Col2", "A", "B",
	}

	for _, w := range want {
		if !contains(result, w) {
			t.Errorf("Render() output should contain %q, got: %q", w, result)
		}
	}
}
