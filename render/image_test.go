package render

import (
	"testing"
)

func TestRender_Images(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "Simple image",
			markdown: "![Alt text](https://example.com/image.png)",
			want:     []string{"Alt text", "https://example.com/image.png", "[Image:"},
		},
		{
			name:     "Image with title",
			markdown: "![Logo](logo.png \"Company Logo\")",
			want:     []string{"Logo", "logo.png"},
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
