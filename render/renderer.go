package render

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

const maxLineWidth = 100

// wrapText wraps text to maxLineWidth characters, breaking at word boundaries when possible
func wrapText(text string) string {
	wrapped, _ := wrapTextWithOffset(text, 0)
	return wrapped
}

// wrapTextWithOffset wraps text to maxLineWidth characters, considering current line offset
func wrapTextWithOffset(text string, currentOffset int) (string, int) {
	if len(text) == 0 {
		return text, currentOffset
	}

	var result strings.Builder
	words := strings.Fields(text)
	currentLine := ""
	lineLength := currentOffset

	// If we're already at or over the limit, start on a new line
	if lineLength >= maxLineWidth && len(words) > 0 {
		result.WriteString("\n")
		lineLength = 0
	}

	for _, word := range words {
		// Handle words that are longer than maxLineWidth by breaking them
		if len(word) > maxLineWidth {
			// Finish current line if it has content
			if currentLine != "" {
				result.WriteString(currentLine)
				result.WriteString("\n")
				currentLine = ""
				lineLength = 0
			}
			// Break the long word into chunks
			for len(word) > maxLineWidth {
				result.WriteString(word[:maxLineWidth])
				result.WriteString("\n")
				word = word[maxLineWidth:]
				lineLength = 0
			}
			if len(word) > 0 {
				currentLine = word
				lineLength = len(word)
			}
			continue
		}

		// If adding this word would exceed the limit, start a new line
		spaceNeeded := 0
		if currentLine != "" {
			spaceNeeded = 1 // space between words
		}
		// Wrap if adding this word would exceed or reach exactly the limit (since maxLineWidth is the maximum)
		if lineLength+len(word)+spaceNeeded >= maxLineWidth {
			if currentLine != "" {
				result.WriteString(currentLine)
				result.WriteString("\n")
				currentLine = word
				lineLength = len(word)
			} else {
				// Current line is empty but we're at/over the limit, wrap to new line
				result.WriteString("\n")
				currentLine = word
				lineLength = len(word)
			}
		} else {
			if currentLine != "" {
				currentLine += " " + word
				lineLength += spaceNeeded + len(word)
			} else {
				currentLine = word
				lineLength += len(word)
			}
		}
	}

	if currentLine != "" {
		result.WriteString(currentLine)
	}

	return result.String(), lineLength
}

// Render renders markdown content with ANSI colors
func Render(content string) {
	// Force color output even when piped (for use with less -R)
	color.NoColor = false

	// Parse markdown
	doc := markdown.Parse([]byte(content), nil)

	// Create renderer
	renderer := &ANSIRenderer{
		listLevel:          0,
		listIndex:          make(map[int]int),
		inCodeBlock:        false,
		inEmph:             false,
		inStrong:           false,
		inHeading:          0,
		currentLineLen:     0,
		justAddedEmphSpace: false,
	}

	// Render and print
	output := renderer.RenderNode(doc)
	fmt.Print(output)
}

// ANSIRenderer renders markdown to ANSI colored terminal output
type ANSIRenderer struct {
	listLevel          int
	listIndex          map[int]int
	inCodeBlock        bool
	inEmph             bool
	inStrong           bool
	inHeading          int  // Track which heading level we're in (0 = not in heading)
	currentLineLen     int  // Track current visual line length (excluding ANSI codes)
	justAddedEmphSpace bool // Track if we just added a space after emphasis
}

// RenderNode recursively renders AST nodes
func (r *ANSIRenderer) RenderNode(node ast.Node) string {
	var buf bytes.Buffer

	ast.WalkFunc(node, func(node ast.Node, entering bool) ast.WalkStatus {
		switch n := node.(type) {
		case *ast.Document:
			// Root node, continue

		case *ast.Heading:
			if entering {
				buf.WriteString("\n")
				r.inHeading = n.Level
				r.currentLineLen = 0
				// Show the # symbols in blue
				prefixLen := 0
				switch n.Level {
				case 1:
					buf.WriteString(color.BlueString("# "))
					prefixLen = 2
				case 2:
					buf.WriteString(color.BlueString("## "))
					prefixLen = 3
				case 3:
					buf.WriteString(color.BlueString("### "))
					prefixLen = 4
				case 4:
					buf.WriteString(color.BlueString("#### "))
					prefixLen = 5
				case 5:
					buf.WriteString(color.BlueString("##### "))
					prefixLen = 6
				case 6:
					buf.WriteString(color.BlueString("###### "))
					prefixLen = 7
				}
				r.currentLineLen = prefixLen
			} else {
				buf.WriteString("\n")
				r.inHeading = 0
				r.currentLineLen = 0
			}

		case *ast.Paragraph:
			if entering {
				r.currentLineLen = 0
			} else {
				buf.WriteString("\n")
				r.currentLineLen = 0
			}

		case *ast.Text:
			if entering {
				text := string(n.Literal)

				// If we just added a space after emphasis and this text starts with a space, skip the leading space to avoid double spaces
				if r.justAddedEmphSpace && len(text) > 0 && text[0] == ' ' {
					text = text[1:]
					r.currentLineLen--
					r.justAddedEmphSpace = false
				} else {
					r.justAddedEmphSpace = false
				}

				wrappedText, newLineLen := wrapTextWithOffset(text, r.currentLineLen)

				// Handle heading text - apply white bold color
				if r.inHeading > 0 {
					buf.WriteString(color.New(color.FgWhite, color.Bold).Sprint(wrappedText))
					// Update line length (count only visible characters, not ANSI codes)
					r.currentLineLen = newLineLen
					return ast.GoToNext
				}

				// Apply formatting based on context for regular text
				if r.inStrong && r.inEmph {
					buf.WriteString(color.New(color.Bold, color.Italic, color.FgHiBlue).Sprint(wrappedText))
				} else if r.inStrong {
					buf.WriteString(color.New(color.Bold, color.FgHiBlue).Sprint(wrappedText))
				} else if r.inEmph {
					buf.WriteString(color.New(color.Italic, color.FgHiBlue).Sprint(wrappedText))
				} else {
					buf.WriteString(wrappedText)
				}

				// Update line length (count only visible characters, not ANSI codes)
				// If wrappedText contains newlines, we're on a new line
				if strings.Contains(wrappedText, "\n") {
					lines := strings.Split(wrappedText, "\n")
					lastLine := lines[len(lines)-1]
					r.currentLineLen = len(lastLine)
				} else {
					r.currentLineLen = newLineLen
				}
			}

		case *ast.Emph:
			if entering {
				// Add space before emphasized text if there's already content on the line
				if r.currentLineLen > 0 {
					buf.WriteString(" ")
					r.currentLineLen++
				}
				r.inEmph = true
				r.justAddedEmphSpace = false
			} else {
				r.inEmph = false
				// Add space after emphasized text
				buf.WriteString(" ")
				r.currentLineLen++
				r.justAddedEmphSpace = true
			}

		case *ast.Strong:
			if entering {
				// Add space before strong text if there's already content on the line
				if r.currentLineLen > 0 {
					buf.WriteString(" ")
					r.currentLineLen++
				}
				r.inStrong = true
				r.justAddedEmphSpace = false
			} else {
				r.inStrong = false
				// Add space after strong text
				buf.WriteString(" ")
				r.currentLineLen++
				r.justAddedEmphSpace = true
			}

		case *ast.Link:
			if entering {
				buf.WriteString(color.BlueString(""))
			} else {
				url := string(n.Destination)
				// Truncate long URLs to fit within maxLineWidth
				urlDisplayLen := len(url)
				if urlDisplayLen > maxLineWidth-10 {
					url = url[:maxLineWidth-10] + "..."
					urlDisplayLen = maxLineWidth - 7
				}
				linkText := fmt.Sprintf(" (%s)", url)
				linkTextLen := len(linkText)

				// Check if adding this link would exceed the line width
				// Wrap if current line + link would exceed, or if we're already at/over the limit
				if r.currentLineLen > 0 {
					if r.currentLineLen+linkTextLen > maxLineWidth || r.currentLineLen >= maxLineWidth {
						buf.WriteString("\n")
						r.currentLineLen = 0
					}
				}

				buf.WriteString(color.New(color.Faint).Sprintf(linkText))
				// Update line length (format: " (url)")
				r.currentLineLen += linkTextLen
				if r.currentLineLen > maxLineWidth {
					// Would exceed, but we already truncated
					r.currentLineLen = maxLineWidth
				}
			}

		case *ast.Image:
			if entering {
				buf.WriteString(color.MagentaString("[Image: "))
				r.currentLineLen += 8 // "[Image: "
			} else {
				url := string(n.Destination)
				// Truncate long image URLs to fit within maxLineWidth
				urlDisplayLen := len(url)
				if urlDisplayLen > maxLineWidth-15 {
					url = url[:maxLineWidth-15] + "..."
					urlDisplayLen = maxLineWidth - 12
				}
				imageText := fmt.Sprintf(" - %s", url)
				buf.WriteString(color.New(color.Faint).Sprintf(imageText))
				buf.WriteString(color.MagentaString("]"))
				// Update line length
				r.currentLineLen += len(imageText) + 1 // +1 for "]"
				if r.currentLineLen > maxLineWidth {
					r.currentLineLen = maxLineWidth
				}
			}

		case *ast.Code:
			if entering {
				code := string(n.Literal)
				// Truncate very long inline code to fit within maxLineWidth
				codeDisplayLen := len(code)
				if codeDisplayLen > maxLineWidth-2 {
					code = code[:maxLineWidth-5] + "..."
					codeDisplayLen = maxLineWidth - 2
				}
				codeText := " " + code + " "
				codeTextLen := len(codeText)

				// Check if adding this code would exceed the line width
				// Wrap if current line + code would exceed, or if we're already at/over the limit
				if r.currentLineLen > 0 {
					if r.currentLineLen+codeTextLen > maxLineWidth || r.currentLineLen >= maxLineWidth {
						buf.WriteString("\n")
						r.currentLineLen = 0
					}
				}

				buf.WriteString(color.New(color.FgHiRed).Sprint(codeText))
				// Update line length
				r.currentLineLen += codeTextLen
				if r.currentLineLen > maxLineWidth {
					r.currentLineLen = maxLineWidth
				}
			}

		case *ast.CodeBlock:
			if entering {
				r.inCodeBlock = true
				boxWidth := maxLineWidth
				buf.WriteString("\n")
				buf.WriteString(color.New(color.FgHiBlack).Sprint("┌" + strings.Repeat("─", boxWidth) + "┐\n"))
				lines := strings.Split(string(n.Literal), "\n")
				for i, line := range lines {
					// Skip the last line if it's empty (trailing newline)
					if i == len(lines)-1 && line == "" {
						continue
					}
					// Wrap long lines within code blocks
					if len(line) > boxWidth-2 {
						// Split long lines
						for len(line) > boxWidth-2 {
							chunk := line[:boxWidth-2]
							line = line[boxWidth-2:]
							buf.WriteString(color.New(color.FgHiBlack).Sprint("│ "))
							buf.WriteString(color.New(color.FgHiMagenta).Sprint(chunk))
							buf.WriteString(color.New(color.FgHiBlack).Sprint(" │\n"))
						}
					}
					// Pad the line to ensure the right border aligns
					paddedLine := line
					if len(line) < boxWidth-2 {
						paddedLine = line + strings.Repeat(" ", boxWidth-2-len(line))
					}
					buf.WriteString(color.New(color.FgHiBlack).Sprint("│ "))
					buf.WriteString(color.New(color.FgHiMagenta).Sprint(paddedLine))
					buf.WriteString(color.New(color.FgHiBlack).Sprint(" │\n"))
				}
				buf.WriteString(color.New(color.FgHiBlack).Sprint("└" + strings.Repeat("─", boxWidth) + "┘\n"))
				r.currentLineLen = 0
			} else {
				r.inCodeBlock = false
			}

		case *ast.List:
			if entering {
				r.listLevel++
				r.listIndex[r.listLevel] = 0
			} else {
				r.listLevel--
				buf.WriteString("\n")
				r.currentLineLen = 0
			}

		case *ast.ListItem:
			if entering {
				r.listIndex[r.listLevel]++
				indent := strings.Repeat("  ", r.listLevel-1)
				indentLen := len(indent)

				// Check if parent is ordered list
				parent := n.GetParent()
				if list, ok := parent.(*ast.List); ok && list.ListFlags&ast.ListTypeOrdered != 0 {
					prefix := fmt.Sprintf("%d. ", r.listIndex[r.listLevel])
					buf.WriteString(indent + color.YellowString(prefix))
					r.currentLineLen = indentLen + len(prefix)
				} else {
					prefix := "• "
					buf.WriteString(indent + color.YellowString(prefix))
					r.currentLineLen = indentLen + len(prefix)
				}
			} else {
				buf.WriteString("\n")
				r.currentLineLen = 0
			}

		case *ast.BlockQuote:
			if entering {
				buf.WriteString(color.New(color.FgHiBlack).Sprint("│ "))
				r.currentLineLen += 2 // "│ "
			}

		case *ast.HorizontalRule:
			if entering {
				buf.WriteString("\n")
				buf.WriteString(color.New(color.FgHiBlack).Sprint(strings.Repeat("─", maxLineWidth)))
				buf.WriteString("\n\n")
				r.currentLineLen = 0
			}

		case *ast.Softbreak, *ast.Hardbreak:
			if entering {
				buf.WriteString("\n")
				r.currentLineLen = 0
			}
		}

		return ast.GoToNext
	})

	return buf.String()
}
