package render

import (
	"testing"
)

func TestRender_Code(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     []string
	}{
		{
			name:     "Inline code",
			markdown: "Use `go build` to build",
			want:     []string{"go build"},
		},
		{
			name:     "Code block",
			markdown: "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
			want:     []string{"func main()", "fmt.Println", "Hello"},
		},
		{
			name:     "Multiple inline code",
			markdown: "Use `var` and `const` for declarations",
			want:     []string{"var", "const"},
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
