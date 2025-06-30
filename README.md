# goflu

A command-line tool to convert Confluence HTML export files to clean Markdown format.

## Installation

```bash
go install github.com/vaihdass/goflu@latest
```

## Usage

Convert a Confluence HTML file to Markdown:

```bash
goflu md path/to/confluence.html
```

Options:
- `-o, --output`: Specify output file path (default: input file with .md extension)
- `-f, --force`: Overwrite output file if it exists

Examples:

```bash
# Convert with default output name
goflu md ./FAQ.html
# Creates ./FAQ.md

# Specify custom output file
goflu md ./FAQ.html -o ./docs/faq.md

# Force overwrite existing file
goflu md ./FAQ.html -f
```

## Features

- Extracts main content from Confluence HTML exports
- Converts common HTML elements to Markdown:
  - Headings (h1-h6)
  - Paragraphs
  - Lists (ordered and unordered)
  - Links
  - Code blocks and inline code
  - Tables
  - Blockquotes
  - Bold and italic text
- Filters out navigation elements and UI components
- Handles both relative and absolute file paths