package render

import (
	"testing"
)

func TestRender_Lists(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "Unordered list",
			markdown: "- Item 1\n- Item 2\n- Item 3",
			want:     []string{"Item 1", "Item 2", "Item 3", "•"},
		},
		{
			name:     "Ordered list",
			markdown: "1. First\n2. Second\n3. Third",
			want:     []string{"First", "Second", "Third", "1.", "2.", "3."},
		},
		{
			name:     "Nested unordered list",
			markdown: "- Item 1\n  - Nested 1\n  - Nested 2\n- Item 2",
			want:     []string{"Item 1", "Nested 1", "Nested 2", "Item 2"},
		},
		{
			name:     "Nested ordered list",
			markdown: "1. First\n   1. Nested A\n   2. Nested B\n2. Second",
			want:     []string{"First", "Nested A", "Nested B", "Second"},
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
