package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

// Render renders markdown content with ANSI colors
func Render(content string) {
	// Parse markdown
	doc := markdown.Parse([]byte(content), nil)

	// Create renderer
	renderer := &ANSIRenderer{
		listLevel:   0,
		listIndex:   make(map[int]int),
		inCodeBlock: false,
		inEmph:      false,
		inStrong:    false,
		inHeading:   0,
	}

	// Render and print
	output := renderer.RenderNode(doc)
	fmt.Print(output)
}

// ANSIRenderer renders markdown to ANSI colored terminal output
type ANSIRenderer struct {
	listLevel   int
	listIndex   map[int]int
	inCodeBlock bool
	inEmph      bool
	inStrong    bool
	inHeading   int // Track which heading level we're in (0 = not in heading)
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
				// Show the # symbols in blue
				switch n.Level {
				case 1:
					buf.WriteString(color.BlueString("# "))
				case 2:
					buf.WriteString(color.BlueString("## "))
				case 3:
					buf.WriteString(color.BlueString("### "))
				case 4:
					buf.WriteString(color.BlueString("#### "))
				case 5:
					buf.WriteString(color.BlueString("##### "))
				case 6:
					buf.WriteString(color.BlueString("###### "))
				}
			} else {
				buf.WriteString("\n")
				r.inHeading = 0
			}

		case *ast.Paragraph:
			if !entering {
				buf.WriteString("\n")
			}

		case *ast.Text:
			if entering {
				text := string(n.Literal)

				// Handle heading text - apply white bold color
				if r.inHeading > 0 {
					buf.WriteString(color.New(color.FgWhite, color.Bold).Sprint(text))
					return ast.GoToNext
				}

				// Apply formatting based on context for regular text
				if r.inStrong && r.inEmph {
					buf.WriteString(color.New(color.Bold, color.FgYellow).Sprint(text))
				} else if r.inStrong {
					buf.WriteString(color.New(color.Bold).Sprint(text))
				} else if r.inEmph {
					buf.WriteString(color.YellowString(text))
				} else {
					buf.WriteString(text)
				}
			}

		case *ast.Emph:
			if entering {
				r.inEmph = true
			} else {
				r.inEmph = false
			}

		case *ast.Strong:
			if entering {
				r.inStrong = true
			} else {
				r.inStrong = false
			}

		case *ast.Link:
			if entering {
				buf.WriteString(color.BlueString(""))
			} else {
				buf.WriteString(color.New(color.Faint).Sprintf(" (%s)", string(n.Destination)))
			}

		case *ast.Image:
			if entering {
				buf.WriteString(color.MagentaString("[Image: "))
			} else {
				buf.WriteString(color.New(color.Faint).Sprintf(" - %s", string(n.Destination)))
				buf.WriteString(color.MagentaString("]"))
			}

		case *ast.Code:
			if entering {
				buf.WriteString(color.GreenString("`" + string(n.Literal) + "`"))
			}

		case *ast.CodeBlock:
			if entering {
				boxWidth := 100
				buf.WriteString("\n")
				buf.WriteString(color.New(color.FgHiBlack).Sprint("┌" + strings.Repeat("─", boxWidth) + "┐\n"))
				lines := strings.Split(string(n.Literal), "\n")
				for i, line := range lines {
					// Skip the last line if it's empty (trailing newline)
					if i == len(lines)-1 && line == "" {
						continue
					}
					// Pad the line to ensure the right border aligns
					paddedLine := line
					if len(line) < boxWidth-2 {
						paddedLine = line + strings.Repeat(" ", boxWidth-2-len(line))
					} else if len(line) > boxWidth-2 {
						paddedLine = line[:boxWidth-2]
					}
					buf.WriteString(color.New(color.FgHiBlack).Sprint("│ "))
					buf.WriteString(color.New(color.FgHiBlack).Sprint(paddedLine + " │\n"))
				}
				buf.WriteString(color.New(color.FgHiBlack).Sprint("└" + strings.Repeat("─", boxWidth) + "┘\n"))
			}

		case *ast.List:
			if entering {
				r.listLevel++
				r.listIndex[r.listLevel] = 0
			} else {
				r.listLevel--
				buf.WriteString("\n")
			}

		case *ast.ListItem:
			if entering {
				r.listIndex[r.listLevel]++
				indent := strings.Repeat("  ", r.listLevel-1)

				// Check if parent is ordered list
				parent := n.GetParent()
				if list, ok := parent.(*ast.List); ok && list.ListFlags&ast.ListTypeOrdered != 0 {
					buf.WriteString(indent + color.YellowString(fmt.Sprintf("%d. ", r.listIndex[r.listLevel])))
				} else {
					buf.WriteString(indent + color.YellowString("• "))
				}
			} else {
				buf.WriteString("\n")
			}

		case *ast.BlockQuote:
			if entering {
				buf.WriteString(color.New(color.FgHiBlack).Sprint("│ "))
			}

		case *ast.HorizontalRule:
			if entering {
				buf.WriteString("\n")
				buf.WriteString(color.New(color.FgHiBlack).Sprint(strings.Repeat("─", 100)))
				buf.WriteString("\n\n")
			}

		case *ast.Softbreak, *ast.Hardbreak:
			if entering {
				buf.WriteString("\n")
			}
		}

		return ast.GoToNext
	})

	return buf.String()
}
