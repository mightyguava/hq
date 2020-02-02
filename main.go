package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/andybalholm/cascadia"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

// Runs a query selector query using goquery on the given file
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cli := struct {
		Selector string `arg:"" help:"the CSS selector to run"`
		In       string `arg:"" help:"file to read input from. If unset, reads from stdin" optional:""`
		Number   int    `help:"Print only the n-th match" short:"n"`
		Text     bool   `help:"Recursively print the inner text instead of HTML"`
	}{}
	kong.Parse(&cli)

	f, err := os.Open(cli.In)
	if err != nil {
		return err
	}
	defer f.Close()
	doc, err := html.Parse(f)
	if err != nil {
		return err
	}
	selector, err := cascadia.Compile(cli.Selector)
	if err != nil {
		return err
	}
	nodes := cascadia.QueryAll(doc, selector)
	buf := bytes.NewBuffer(nil)
	if cli.Number != 0 && cli.Number > len(nodes) {
		return fmt.Errorf("wanted match %d, but only got %d matches", cli.Number, len(nodes))
	}
	if cli.Number > 0 {
		nodes = []*html.Node{nodes[cli.Number-1]}
	}
	if cli.Text {
		for _, n := range nodes {
			buf := &bytes.Buffer{}
			walk(n, func(node *html.Node) {
				if node.Type == html.TextNode {
					buf.WriteString(node.Data)
				}
			})
			fmt.Println(buf.String())
		}
	} else {
		for _, n := range nodes {
			if err := html.Render(buf, n); err != nil {
				return err
			}
		}
		fmt.Println(string(gohtml.FormatBytes(buf.Bytes())))
	}
	return nil
}

func walk(n *html.Node, visit func(node *html.Node)) {
	visit(n)
	if n.FirstChild == nil {
		return
	}
	for n = n.FirstChild; n != nil; n = n.NextSibling {
		walk(n, visit)
	}
}
