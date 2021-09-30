package rel2abs

import (
	"bytes"
	"fmt"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func rel2abs(n *html.Node, nurl *url.URL) error {
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		for i := range n.Attr {
			if n.Attr[i].Key == "href" {
				rel, err := url.Parse(n.Attr[i].Val)
				if err != nil {
					return fmt.Errorf("relative url: %w\n", err)
				}

				n.Attr[i].Val = nurl.ResolveReference(rel).String()
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		rel2abs(c, nurl)
	}
	return nil
}

// Converts all relative URLs in htmlContent to absolute URLs,
// resolved against a base URL.
func Convert(htmlContent []byte, base string) ([]byte, error) {
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
