# mdrender

A simple, beautiful markdown renderer for your terminal.

## Features

- 🎨 Beautiful ANSI color output
- 📁 Read from files
- 🌐 Fetch and render from URLs
- ⌨️  Read from stdin
- 🚀 Fast and lightweight
- 📝 Clean, readable code

## Installation

```bash
go install github.com/yourusername/mdrender@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/mdrender.git
cd mdrender
go build
```

## Usage

### Render a markdown file

```bash
mdrender README.md
mdrender path/to/file.md
```

### Render from URL

```bash
mdrender https://raw.githubusercontent.com/user/repo/main/README.md
```

### Render from stdin

```bash
echo "# Hello World" | mdrender
cat file.md | mdrender
```

### Render inline markdown

```bash
mdrender "# Hello\nThis is **bold** text"
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

## Examples

See `examples/demo.md` for a comprehensive example.

## License

MIT License - feel free to use this in your daily workflow!

## Contributing

Contributions welcome! Keep it simple and maintainable.
