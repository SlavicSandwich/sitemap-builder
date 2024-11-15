package parser

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func extractHref(cur *html.Node) (*string, error) {
	for _, attr := range cur.Attr {
		if attr.Key == "href" {
			return &attr.Val, nil
		}
	}
	return nil, fmt.Errorf("couldn't find 'href'")
}

func extractText(cur *html.Node) (*string, error) {
	var res []string
	for c := cur.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			res = append(res, c.Data)
		}
	}
	str := strings.TrimSpace(strings.Join(res, ""))
	return &str, nil
}

func newLink(node *html.Node) (*Link, error) {
	href, err := extractHref(node)
	if err != nil {
		return nil, err
	}
	text, err := extractText(node)
	if err != nil {
		return nil, err
	}
	return &Link{*href, *text}, nil
}

func Parse(r io.Reader) ([]Link, error) {
	var arr []Link
	head, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	err = dfs(head, &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func dfs(cur *html.Node, arr *[]Link) error {
	if cur == nil {
		return nil
	}
	if cur.Type == html.ElementNode && cur.Data == "a" {
		link, err := newLink(cur)
		if err != nil {
			return err
		}
		*arr = append(*arr, *link)
		return nil
	}

	for c := cur.FirstChild; c != nil; c = c.NextSibling {
		err := dfs(c, arr)
		if err != nil {
			return err
		}
	}
	return nil
}
