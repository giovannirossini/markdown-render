package render

import (
	"testing"
)

func TestRender_HorizontalRule(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		{
			name:     "Horizontal rule",
			markdown: "---",
			want:     "─",
		},
		{
			name:     "Horizontal rule with content",
			markdown: "Above\n\n---\n\nBelow",
			want:     "─",
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
