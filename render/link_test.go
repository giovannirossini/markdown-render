package render

import (
	"testing"
)

func TestRender_Links(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "Simple link",
			markdown: "[Google](https://google.com)",
			want:     []string{"Google", "https://google.com"},
		},
		{
			name:     "Link with text",
			markdown: "Visit [GitHub](https://github.com) for code",
			want:     []string{"GitHub", "https://github.com"},
		},
		{
			name:     "Multiple links",
			markdown: "[Link1](http://example.com) and [Link2](http://test.com)",
			want:     []string{"Link1", "http://example.com", "Link2", "http://test.com"},
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
