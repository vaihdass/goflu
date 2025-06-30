package parser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseConfluenceHTML converts Confluence HTML export to Markdown
func ParseConfluenceHTML(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var markdown strings.Builder

	// Extract title
	title := doc.Find("title").Text()
	if title != "" {
		markdown.WriteString("# " + cleanText(title) + "\n\n")
	}

	// Look for main content area
	// Confluence usually has content in div with id="main-content" or class="wiki-content"
	mainContent := doc.Find("#main-content, .wiki-content").First()

	if mainContent.Length() == 0 {
		// Fallback: try to find any content area
		mainContent = doc.Find("div[class*='content']").First()
	}

	if mainContent.Length() > 0 {
		parseChildren(mainContent, &markdown, 0)
	} else {
		// Last resort: parse body
		parseChildren(doc.Find("body"), &markdown, 0)
	}

	return strings.TrimSpace(markdown.String()), nil
}

func parseChildren(s *goquery.Selection, markdown *strings.Builder, listLevel int) {
	s.Children().Each(func(i int, node *goquery.Selection) {
		parseElement(node, markdown, listLevel)
	})
}

func parseElement(s *goquery.Selection, markdown *strings.Builder, listLevel int) {
	if s.Is("h1, h2, h3, h4, h5, h6") {
		level := s.Get(0).Data[1] - '0'
		text := cleanText(s.Text())
		markdown.WriteString(strings.Repeat("#", int(level)) + " " + text + "\n\n")
	} else if s.Is("p") {
		text := parseInlineElements(s)
		if text != "" {
			markdown.WriteString(text + "\n\n")
		}
	} else if s.Is("ul, ol") {
		parseList(s, markdown, listLevel, s.Is("ol"))
	} else if s.Is("li") {
		// Skip li elements at root level (they should be handled by parseList)
	} else if s.Is("pre") {
		code := s.Find("code").Text()
		if code == "" {
			code = s.Text()
		}
		markdown.WriteString("```\n" + code + "\n```\n\n")
	} else if s.Is("code") && !s.Parent().Is("pre") {
		markdown.WriteString("`" + s.Text() + "`")
	} else if s.Is("blockquote") {
		lines := strings.Split(strings.TrimSpace(parseInlineElements(s)), "\n")
		for _, line := range lines {
			markdown.WriteString("> " + line + "\n")
		}
		markdown.WriteString("\n")
	} else if s.Is("table") {
		parseTable(s, markdown)
	} else if s.Is("hr") {
		markdown.WriteString("---\n\n")
	} else if s.Is("div, section, article") {
		// Check if it's a panel or note
		class, _ := s.Attr("class")
		if strings.Contains(class, "panel") || strings.Contains(class, "note") || strings.Contains(class, "info") {
			markdown.WriteString("> **" + extractPanelTitle(s) + "**: ")
			markdown.WriteString(cleanText(s.Text()) + "\n\n")
		} else {
			parseChildren(s, markdown, listLevel)
		}
	} else if s.Is("a") {
		href, _ := s.Attr("href")
		text := s.Text()
		if text != "" && href != "" {
			markdown.WriteString("[" + text + "](" + href + ")")
		} else if text != "" {
			markdown.WriteString(text)
		}
	} else if s.Is("strong, b") {
		markdown.WriteString("**" + s.Text() + "**")
	} else if s.Is("em, i") {
		markdown.WriteString("*" + s.Text() + "*")
	} else if s.Is("br") {
		markdown.WriteString("\n")
	} else {
		// For other elements, recurse into children
		parseChildren(s, markdown, listLevel)
	}
}

func parseInlineElements(s *goquery.Selection) string {
	var result strings.Builder

	s.Contents().Each(func(i int, node *goquery.Selection) {
		if goquery.NodeName(node) == "#text" {
			result.WriteString(node.Text())
		} else if node.Is("a") {
			href, _ := node.Attr("href")
			text := node.Text()
			if text != "" && href != "" {
				result.WriteString("[" + text + "](" + href + ")")
			} else if text != "" {
				result.WriteString(text)
			}
		} else if node.Is("strong, b") {
			result.WriteString("**" + node.Text() + "**")
		} else if node.Is("em, i") {
			result.WriteString("*" + node.Text() + "*")
		} else if node.Is("code") {
			result.WriteString("`" + node.Text() + "`")
		} else if node.Is("br") {
			result.WriteString("\n")
		} else {
			result.WriteString(parseInlineElements(node))
		}
	})

	return strings.TrimSpace(result.String())
}

func parseList(s *goquery.Selection, markdown *strings.Builder, level int, ordered bool) {
	s.Children().Filter("li").Each(func(i int, item *goquery.Selection) {
		prefix := strings.Repeat("  ", level)
		if ordered {
			prefix += fmt.Sprintf("%d. ", i+1)
		} else {
			prefix += "- "
		}

		// Get direct text content
		var itemText strings.Builder
		item.Contents().Each(func(j int, node *goquery.Selection) {
			if node.Is("ul, ol") {
				// Skip nested lists, they'll be handled separately
			} else if goquery.NodeName(node) == "#text" {
				itemText.WriteString(node.Text())
			} else {
				itemText.WriteString(parseInlineElements(node))
			}
		})

		text := strings.TrimSpace(itemText.String())
		if text != "" {
			markdown.WriteString(prefix + text + "\n")
		}

		// Handle nested lists
		item.Children().Filter("ul, ol").Each(func(j int, nested *goquery.Selection) {
			parseList(nested, markdown, level+1, nested.Is("ol"))
		})
	})

	if level == 0 {
		markdown.WriteString("\n")
	}
}

func parseTable(s *goquery.Selection, markdown *strings.Builder) {
	// Find all rows
	rows := s.Find("tr")
	if rows.Length() == 0 {
		return
	}

	// Process header row (if exists)
	headerCells := rows.First().Find("th")
	if headerCells.Length() > 0 {
		markdown.WriteString("|")
		headerCells.Each(func(i int, cell *goquery.Selection) {
			markdown.WriteString(" " + cleanText(cell.Text()) + " |")
		})
		markdown.WriteString("\n|")
		headerCells.Each(func(i int, cell *goquery.Selection) {
			markdown.WriteString(" --- |")
		})
		markdown.WriteString("\n")
		rows = rows.Slice(1, rows.Length())
	}

	// Process data rows
	rows.Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() > 0 {
			markdown.WriteString("|")
			cells.Each(func(j int, cell *goquery.Selection) {
				markdown.WriteString(" " + cleanText(cell.Text()) + " |")
			})
			markdown.WriteString("\n")
		}
	})

	markdown.WriteString("\n")
}

func extractPanelTitle(s *goquery.Selection) string {
	// Try to find panel title
	title := s.Find(".panelHeader, .panel-heading, .title").First().Text()
	if title == "" {
		class, _ := s.Attr("class")
		if strings.Contains(class, "info") {
			return "Info"
		} else if strings.Contains(class, "warning") {
			return "Warning"
		} else if strings.Contains(class, "note") {
			return "Note"
		}
		return "Note"
	}
	return cleanText(title)
}

func cleanText(text string) string {
	// Remove excessive whitespace
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Replace multiple spaces with single space
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}
