package render

import (
	"testing"
)

func TestRender_Headings(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []string // substrings that should be present
	}{
		{
			name:     "H1 heading",
			markdown: "# Hello World",
			want:     []string{"# ", "Hello World"},
		},
		{
			name:     "H2 heading",
			markdown: "## Second Level",
			want:     []string{"## ", "Second Level"},
		},
		{
			name:     "H3 heading",
			markdown: "### Third Level",
			want:     []string{"### ", "Third Level"},
		},
		{
			name:     "H4 heading",
			markdown: "#### Fourth Level",
			want:     []string{"#### ", "Fourth Level"},
		},
		{
			name:     "H5 heading",
			markdown: "##### Fifth Level",
			want:     []string{"##### ", "Fifth Level"},
		},
		{
			name:     "H6 heading",
			markdown: "###### Sixth Level",
			want:     []string{"###### ", "Sixth Level"},
		},
		{
			name:     "Multiple headings",
			markdown: "# Title\n## Subtitle\n### Section",
			want:     []string{"# ", "Title", "## ", "Subtitle", "### ", "Section"},
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
