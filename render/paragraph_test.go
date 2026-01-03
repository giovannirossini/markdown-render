package render

import (
	"testing"
)

func TestRender_Paragraphs(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		{
			name:     "Simple paragraph",
			markdown: "This is a paragraph.",
			want:     "This is a paragraph.",
		},
		{
			name:     "Multiple paragraphs",
			markdown: "First paragraph.\n\nSecond paragraph.",
			want:     "First paragraph.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderToString(tt.markdown)

			if !contains(result, tt.want) {
				t.Errorf("RenderToString() output should contain %q, got: %q", tt.want, result)
			}
		})
	}
}
