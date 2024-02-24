package scrape

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Extractor struct {
	IgnoredTags map[string]struct{}
}

func (e *Extractor) ExtractFromReader(r io.Reader) string {
	t := html.NewTokenizer(r)

	var currentToken html.Token
	var result strings.Builder

	for {
		tokenType := t.Next()
		if tokenType == html.ErrorToken {
			break
		}
		if tokenType == html.StartTagToken {
			currentToken = t.Token()
			continue
		}
		if tokenType == html.TextToken {
			_, ok := e.IgnoredTags[currentToken.Data]
			if !ok {
				text := strings.TrimSpace(html.UnescapeString(string(t.Text())))
				if len(text) > 0 {
					result.WriteString(text)
					result.WriteRune('\n')
				}
			}
		}
	}
	return result.String()
}
