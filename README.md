# markdown-render

[![main](https://github.com/giovannirossini/markdown-render/actions/workflows/main.yml/badge.svg)](https://github.com/giovannirossini/markdown-render/actions/workflows/ci.yml) [![release](https://github.com/giovannirossini/markdown-render/actions/workflows/release.yml/badge.svg)](https://github.com/giovannirossini/markdown-render/actions/workflows/ci.yml)

A simple, beautiful markdown renderer for your terminal.

## Features

- 🎨 Beautiful ANSI color output
- 📁 Read from files
- ⌨️  Read from stdin
- 🚀 Fast and lightweight

## Installation

```bash
go install github.com/giovannirossini/markdown-render@latest
```

Or build from source:

```bash
git clone https://github.com/giovannirossini/markdown-render.git
cd markdown-render
make build
```

## Usage

### Render a markdown file

```bash
markdown-render README.md
markdown-render path/to/file.md
```

### Render from stdin

```bash
echo "# Hello World" | markdown-render
cat file.md | markdown-render
```

### Render inline markdown

```bash
markdown-render "# Hello\nThis is **bold** text"
```

## Supported Markdown Features

- ✅ Headings (H1-H6)
- ✅ Bold and italic text
- ✅ Links
- ✅ Images
- ✅ Code blocks and inline code
- ✅ Lists (ordered and unordered)
- ✅ Nested lists
- ✅ Blockquotes
- ✅ Horizontal rules
- ✅ Line breaks

## Color Scheme

- **Headings**: Cyan + Bold
- **Bold text**: Bold
- **Italic text**: Yellow
- **Links**: Blue
- **Code**: Green
- **Images**: Magenta
- **List bullets**: Yellow
- **Blockquotes**: Dark gray
- **Horizontal rules**: Dark gray

