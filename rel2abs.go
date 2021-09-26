package rel2abs

import (
	"bytes"
	"fmt"
	"net/url"

	"golang.org/x/net/html"
)

func rel2abs(n *html.Node, nurl *url.URL) error {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				rel, err := url.Parse(a.Val)
				if err != nil {
					return fmt.Errorf("relative url: %w\n", err)
				}

				a.Val = nurl.ResolveReference(rel).String()
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rel2abs(c, nurl)
		}
	}
	return nil
}

// Converts all relative URLs in htmlContent to absolute URLs,
// resolved against a base URL.
// Example, with base as http://example.com/foo:
//    <a href="#fn-1">
// becomes
//    <a href="http://example.com/foo#fn-1">
func Rel2Abs(htmlContent []byte, base string) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("html parse: %w\n", err)
	}

	nurl, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("url parse: %w\n", err)
	}
	rel2abs(doc, nurl)
	buf := bytes.Buffer{}
	html.Render(&buf, doc)
	return buf.Bytes(), nil
}
