package convert

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"good.site/pkg/utils/queue"
	"good.site/pkg/utils/stack"
)

func bfs(node *html.Node, mark func(node *html.Node) bool) {
	q := queue.New()

	q.Enqueue(node)

	for q.Len() > 0 {
		n := q.Dequeue().(*html.Node)

		done := mark(n)

		if done {
			break
		}

		for n := range n.ChildNodes() {
			q.Enqueue(n)
		}
	}
}

func dfs(node *html.Node, mark func(node *html.Node) bool) {
	st := stack.New()

	st.Push(node)

	for st.Len() > 0 {
		n := st.Pop().(*html.Node)

		done := mark(n)

		if done {
			break
		}

		for n := range n.ChildNodes() {
			st.Push(n)
		}
	}
}

func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	err := html.Render(w, n)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func RemoveNodes(node *html.Node, nodeStr string, count int) []*html.Node {
	removed := []*html.Node{}
	dfs(node, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == nodeStr {
			n.Parent.RemoveChild(n)
			removed = append(removed, n)
			count = count - 1
			if count == 0 {
				return true
			}
		}

		return false
	})

	return removed
}

func FindNodes(node *html.Node, nodeStr string, count int) []*html.Node {
	nodes := []*html.Node{}

	dfs(node, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == nodeStr {
			nodes = append(nodes, n)
			count = count - 1

			if count == 0 {
				return true
			}
		}

		return false
	})

	return nodes
}

func AddStyle(node *html.Node, style string) {
	bfs(node, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "head" {
			n.AppendChild(&html.Node{
				Type: html.ElementNode,
				Data: "style",
				FirstChild: &html.Node{
					Type: html.TextNode,
					Data: style,
				},
			})

			return true
		}

		return false
	})
}

func GetTextContent(node *html.Node) string {
	textContent := ""

	dfs(node, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			textContent = fmt.Sprintf("%s%s", textContent, n.Data)
		}

		return false
	})

	return textContent
}

func TextContentFromHtml(node *html.Node, nodeStr string) string {
	text := ""

	dfs(node, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == nodeStr && len(n.FirstChild.Data) > 0 {
			text = GetTextContent(n)
			return true
		}

		return false
	})
	return text
}

func ExtractStyles(node *html.Node) string {
	style := ""
	styleNodes := RemoveNodes(node, "style", -1)

	for _, n := range styleNodes {
		style = fmt.Sprintf("%s\n%s", style, GetTextContent(n))
	}

	return style
}

func ParseHtml(domStr string) (*html.Node, error) {
	return html.Parse(strings.NewReader(domStr))
}

func ReplaceText(node *html.Node, getText func(input string) string) {
	dfs(node, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			n.Data = getText(n.Data)
		}

		return false
	})
}

func tokenize(node *html.Node) []string {
	tokens := []string{}

	dfs(node, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			tokens = append(tokens, strings.Split(n.Data, " ")...)
		}

		return false
	})

	output := []string{}
	r, err := regexp.Compile(`^[\w ]*[^\W_][\w ]*$`)

	if err != nil {
		panic(err)
	}

	for _, token := range tokens {
		token = strings.TrimSpace(token)

		if _, err = strconv.Atoi(token); err == nil {
			continue
		} else if !r.MatchString(token) {
			continue
		}

		output = append(output, token)
	}

	return output
}

func TokenizeHtml(domStr string) []string {
	doc, err := ParseHtml(domStr)

	if err != nil {
		panic(err)
	}

	body := FindNodes(doc, "body", 1)

	if len(body) > 0 {
		tokens := tokenize(body[0])

		return tokens
	}

	return []string{}
}
