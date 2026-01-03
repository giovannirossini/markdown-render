package render

import (
	"testing"
)

func TestRender_BoldAndItalic(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "Bold text",
			markdown: "This is **bold** text",
			want:     []string{"bold"},
		},
		{
			name:     "Italic text",
			markdown: "This is *italic* text",
			want:     []string{"italic"},
		},
		{
			name:     "Bold and italic",
			markdown: "This is ***bold and italic*** text",
			want:     []string{"bold and italic"},
		},
		{
			name:     "Mixed formatting",
			markdown: "**Bold** and *italic* in one line",
			want:     []string{"Bold", "italic"},
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
