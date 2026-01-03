package render

import (
	"testing"
)

func TestRender_Tables(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "Simple table",
			markdown: "| Header 1 | Header 2 |\n|----------|----------|\n| Cell 1   | Cell 2   |",
			want:     []string{"Header 1", "Header 2", "Cell 1", "Cell 2"},
		},
		{
			name:     "Table with alignment",
			markdown: "| Left | Center | Right |\n|:-----|:------:|------:|\n| A    | B      | C     |",
			want:     []string{"Left", "Center", "Right", "A", "B", "C"},
		},
		{
			name:     "Table with multiple rows",
			markdown: "| Col1 | Col2 |\n|------|------|\n| A    | B    |\n| C    | D    |",
			want:     []string{"Col1", "Col2", "A", "B", "C", "D"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderToString(tt.markdown)

			for _, want := range tt.want {
				if !contains(result, want) {
					t.Errorf("RenderToString() output should contain %q, got: %q", want, result)
				}
			}
		})
	}
}
