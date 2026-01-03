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

// calculateTableColumnWidths calculates the width needed for each column
func (r *ANSIRenderer) calculateTableColumnWidths() {
	if len(r.tableRows) == 0 {
		return
	}

	numCols := len(r.tableRows[0])
	r.tableColumnWidths = make([]int, numCols)

	// Find max width for each column
	for _, row := range r.tableRows {
		for i, cell := range row {
			if i < len(r.tableColumnWidths) {
				cellLen := len(cell)
				if cellLen > r.tableColumnWidths[i] {
					r.tableColumnWidths[i] = cellLen
				}
			}
		}
	}

	// Ensure minimum width of 3 for each column
	for i := range r.tableColumnWidths {
		if r.tableColumnWidths[i] < 3 {
			r.tableColumnWidths[i] = 3
		}
	}
}

// renderTable renders the collected table data
func (r *ANSIRenderer) renderTable() string {
	if len(r.tableRows) == 0 {
		return ""
	}

	var result strings.Builder
	r.calculateTableColumnWidths()

	// Calculate total table width
	totalWidth := 1 // Start with left border
	for _, width := range r.tableColumnWidths {
		totalWidth += width + 3 // cell width + 2 spaces padding + 1 border
	}

	// Limit table width to maxLineWidth
	if totalWidth > maxLineWidth {
		// Scale down columns proportionally
		scale := float64(maxLineWidth-1-len(r.tableColumnWidths)*3) / float64(totalWidth-1-len(r.tableColumnWidths)*3)
		for i := range r.tableColumnWidths {
			r.tableColumnWidths[i] = int(float64(r.tableColumnWidths[i]) * scale)
			if r.tableColumnWidths[i] < 3 {
				r.tableColumnWidths[i] = 3
			}
		}
	}

	// Render top border
	result.WriteString("\n")
	result.WriteString(color.New(color.FgHiBlack).Sprint("┌"))
	for i, width := range r.tableColumnWidths {
		if i > 0 {
			result.WriteString(color.New(color.FgHiBlack).Sprint("┬"))
		}
		result.WriteString(color.New(color.FgHiBlack).Sprint(strings.Repeat("─", width+2)))
	}
	result.WriteString(color.New(color.FgHiBlack).Sprint("┐\n"))

	// Render rows
	for rowIdx, row := range r.tableRows {
		// Render cell row
		result.WriteString(color.New(color.FgHiBlack).Sprint("│"))
		for colIdx, cell := range row {
			cellContent := cell
			if colIdx >= len(r.tableColumnWidths) {
				continue
			}
			colWidth := r.tableColumnWidths[colIdx]

			// Truncate cell content if too long
			if len(cellContent) > colWidth {
				cellContent = cellContent[:colWidth-3] + "..."
			}

			// Apply alignment
			alignment := ast.TableAlignmentLeft
			if colIdx < len(r.tableAlignments) {
				alignment = r.tableAlignments[colIdx]
			}

			var paddedCell string
			switch alignment {
			case ast.TableAlignmentCenter:
				padding := colWidth - len(cellContent)
				leftPad := padding / 2
				rightPad := padding - leftPad
				paddedCell = strings.Repeat(" ", leftPad) + cellContent + strings.Repeat(" ", rightPad)
			case ast.TableAlignmentRight:
				paddedCell = fmt.Sprintf("%*s", colWidth, cellContent)
			default: // Left alignment
				paddedCell = fmt.Sprintf("%-*s", colWidth, cellContent)
			}

			// Apply header styling for first row
			if rowIdx == 0 {
				result.WriteString(" ")
				result.WriteString(color.New(color.FgWhite, color.Bold).Sprint(paddedCell))
				result.WriteString(" ")
			} else {
				result.WriteString(" ")
				result.WriteString(paddedCell)
				result.WriteString(" ")
			}
			result.WriteString(color.New(color.FgHiBlack).Sprint("│"))
		}
		result.WriteString("\n")

		// Render separator after header
		if rowIdx == 0 {
			result.WriteString(color.New(color.FgHiBlack).Sprint("├"))
			for i, width := range r.tableColumnWidths {
				if i > 0 {
					result.WriteString(color.New(color.FgHiBlack).Sprint("┼"))
				}
				result.WriteString(color.New(color.FgHiBlack).Sprint(strings.Repeat("─", width+2)))
			}
			result.WriteString(color.New(color.FgHiBlack).Sprint("┤\n"))
		}
	}

	// Render bottom border
	result.WriteString(color.New(color.FgHiBlack).Sprint("└"))
	for i, width := range r.tableColumnWidths {
		if i > 0 {
			result.WriteString(color.New(color.FgHiBlack).Sprint("┴"))
		}
		result.WriteString(color.New(color.FgHiBlack).Sprint(strings.Repeat("─", width+2)))
	}
	result.WriteString(color.New(color.FgHiBlack).Sprint("┘\n"))

	return result.String()
}

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

// RenderToString renders markdown content with ANSI colors and returns the string
func RenderToString(content string) string {
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
		inTable:            false,
		tableColumnWidths:  nil,
		tableCurrentRow:    nil,
		tableRows:          nil,
		tableAlignments:    nil,
		isTableHeader:      false,
		inTableCell:        false,
		tableCellBuffer:    nil,
	}

	// Render and return
	return renderer.RenderNode(doc)
}

// Render renders markdown content with ANSI colors and prints to stdout
func Render(content string) {
	output := RenderToString(content)
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
	// Table rendering state
	inTable           bool
	tableColumnWidths []int
	tableCurrentRow   []string
	tableRows         [][]string
	tableAlignments   []ast.CellAlignFlags
	isTableHeader     bool
	inTableCell       bool
	tableCellBuffer   *strings.Builder
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

				// If we're in a table cell, collect the text instead of rendering it
				if r.inTableCell && r.tableCellBuffer != nil {
					r.tableCellBuffer.WriteString(text)
					return ast.GoToNext
				}

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
				// If we're in a table cell, just mark emphasis but don't add spaces
				if r.inTableCell {
					r.inEmph = true
				} else {
					// Add space before emphasized text if there's already content on the line
					if r.currentLineLen > 0 {
						buf.WriteString(" ")
						r.currentLineLen++
					}
					r.inEmph = true
					r.justAddedEmphSpace = false
				}
			} else {
				// If we're in a table cell, just unmark emphasis
				if r.inTableCell {
					r.inEmph = false
				} else {
					r.inEmph = false
					// Add space after emphasized text
					buf.WriteString(" ")
					r.currentLineLen++
					r.justAddedEmphSpace = true
				}
			}

		case *ast.Strong:
			if entering {
				// If we're in a table cell, just mark strong but don't add spaces
				if r.inTableCell {
					r.inStrong = true
				} else {
					// Add space before strong text if there's already content on the line
					if r.currentLineLen > 0 {
						buf.WriteString(" ")
						r.currentLineLen++
					}
					r.inStrong = true
					r.justAddedEmphSpace = false
				}
			} else {
				// If we're in a table cell, just unmark strong
				if r.inTableCell {
					r.inStrong = false
				} else {
					r.inStrong = false
					// Add space after strong text
					buf.WriteString(" ")
					r.currentLineLen++
					r.justAddedEmphSpace = true
				}
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

				// If we're in a table cell, collect the code instead of rendering it
				if r.inTableCell && r.tableCellBuffer != nil {
					r.tableCellBuffer.WriteString(code)
					return ast.GoToNext
				}

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

		case *ast.Table:
			if entering {
				r.inTable = true
				r.tableRows = make([][]string, 0)
				r.tableCurrentRow = nil
				r.tableColumnWidths = nil
				r.tableAlignments = nil
			} else {
				// Render the complete table
				tableOutput := r.renderTable()
				buf.WriteString(tableOutput)
				r.inTable = false
				r.tableRows = nil
				r.tableCurrentRow = nil
				r.tableColumnWidths = nil
				r.tableAlignments = nil
				r.currentLineLen = 0
			}

		case *ast.TableHeader:
			if entering {
				r.isTableHeader = true
			} else {
				r.isTableHeader = false
			}

		case *ast.TableBody:
			// TableBody is just a container, no special handling needed
			// entering/leaving doesn't need special logic

		case *ast.TableRow:
			if entering {
				r.tableCurrentRow = make([]string, 0)
			} else {
				// Row complete, add it to tableRows
				if len(r.tableCurrentRow) > 0 {
					r.tableRows = append(r.tableRows, r.tableCurrentRow)
				}
				r.tableCurrentRow = nil
			}

		case *ast.TableCell:
			if entering {
				// Start collecting cell content
				r.inTableCell = true
				r.tableCellBuffer = &strings.Builder{}
				// Store alignment for first row (header)
				if r.isTableHeader && len(r.tableAlignments) < len(r.tableCurrentRow)+1 {
					r.tableAlignments = append(r.tableAlignments, n.Align)
				}
			} else {
				// Cell complete, add content to current row
				cellContent := strings.TrimSpace(r.tableCellBuffer.String())
				r.tableCurrentRow = append(r.tableCurrentRow, cellContent)
				r.inTableCell = false
				r.tableCellBuffer = nil
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
