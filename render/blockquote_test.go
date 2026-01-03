package render

import (
	"testing"
)

func TestRender_Blockquotes(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "Simple blockquote",
			markdown: "> This is a quote",
			want:     []string{"This is a quote"},
		},
		{
			name:     "Multi-line blockquote",
			markdown: "> Line 1\n> Line 2",
			want:     []string{"Line 1", "Line 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderToString(tt.markdown)

			for _, want := range tt.want {
				if !contains(result, want) {
					t.Errorf("RenderToString() output should contain %q, got: %q", tt.want, result)
				}
			}
		})
	}
}
