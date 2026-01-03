package render

import (
	"testing"
)

func TestRender_EmptyInput(t *testing.T) {
	result := RenderToString("")

	// Empty input should produce minimal output (just newlines from formatting)
	if len(result) > 10 {
		t.Errorf("RenderToString() with empty input should produce minimal output, got: %q", result)
	}
}
