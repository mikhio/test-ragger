package htmlx

import (
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var cleanReSpace = regexp.MustCompile(`\s+`)

// ToText parses an HTML file and returns normalized text and title.
func ToText(r io.Reader, fallbackPath string) (string, string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", "", err
	}

	// remove noise
	doc.Find("script,style,nav,footer,header,noscript").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	title := strings.TrimSpace(doc.Find("title").First().Text())
	if title == "" {
		title = filepath.Base(fallbackPath)
	}

	text := strings.TrimSpace(doc.Text())
	text = cleanReSpace.ReplaceAllString(text, " ")
	return text, title, nil
}
