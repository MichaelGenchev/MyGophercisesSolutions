package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	inputHTML, err := os.ReadFile("ex4.html")
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(bytes.NewReader(inputHTML))
	if err != nil {
		log.Fatal(err)
	}

	links, err := extractLinks(doc)
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range links {
		fmt.Println(l)
	}
}

type Link struct {
	Href string
	Text string
}

func extractLinks(doc *html.Node) ([]Link, error) {
	// Start DFS to group the links
	var hrefs []Link

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node == nil {
			return
		}

		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					hrefs = append(hrefs, Link{Href: attr.Val, Text: collectText(node)})
					break
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)

	return hrefs, nil
}

func collectText(node *html.Node) string {
	var b strings.Builder
	collectInto(node, &b)
	return b.String()
}

func collectInto(n *html.Node, b *strings.Builder) {
	if n == nil {
		return
	}

	if n.Type == html.TextNode {
		b.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectInto(c, b)
	}
}
