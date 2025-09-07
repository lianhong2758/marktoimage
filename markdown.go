package marktoimage

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

func NewRichTextFromMarkdown(content string) *RichText {
	rt := NewRichText(
		Config{
			Width:      1080,
			Height:     1920,
			LineHeight: 8,
			TextMargin: 20,
			TopMargin:  10,
		})
	rt.AppendSegment(parseMarkdown(content)...)
	return rt
}
func (t *RichText) ParseMarkdown(content string) {
	t.Segments = []RichTextSegment{}
	t.AppendMarkdown(content)
}
func (t *RichText) AppendMarkdown(content string) {
	t.AppendSegment(parseMarkdown(content)...)
}

func parseMarkdown(content string) []RichTextSegment {
	r := markdownRenderer{}
	md := goldmark.New(goldmark.WithRenderer(&r))
	err := md.Convert([]byte(content), nil)
	if err != nil {
		log.Print(err)
	}
	return r
}

type markdownRenderer []RichTextSegment

func (m *markdownRenderer) AddOptions(...renderer.Option) {}

func (m *markdownRenderer) Render(_ io.Writer, source []byte, n ast.Node) error {
	segs, err := renderNode(source, n, false)
	*m = segs
	return err
}

func renderNode(source []byte, n ast.Node, blockquote bool) ([]RichTextSegment, error) {
	switch t := n.(type) {
	case *ast.Document:
		return renderChildren(source, n, blockquote)
	case *ast.Paragraph:
		children, err := renderChildren(source, n, blockquote)
		if !blockquote {
			linebreak := &TextSegment{Style: TextStyleParagraph, Text: ""}
			if children[len(children)-1].Inline() {
				if t, ok := children[len(children)-1].(*TextSegment); ok {
					t.Style.Inline = false
				}
			}
			children = append(children, linebreak)
		}
		return children, err
	case *ast.List:
		items, err := renderChildren(source, n, blockquote)
		return []RichTextSegment{
			&ListSegment{Items: items, Ordered: t.Marker != '*' && t.Marker != '-' && t.Marker != '+'},
		}, err
	case *ast.ListItem:
		texts, err := renderChildren(source, n, blockquote)
		return texts, err
	case *ast.TextBlock:
		return renderChildren(source, n, blockquote)
	case *ast.Heading:
		text := forceIntoHeadingText(source, n)
		switch t.Level {
		case 1:
			return []RichTextSegment{&TextSegment{Style: TextStyleHead, Text: text}}, nil
		case 2:
			return []RichTextSegment{&TextSegment{Style: TextStyleSubHead, Text: text}}, nil
		default:
			textSegment := TextSegment{Style: TextStyleParagraph, Text: text}
			textSegment.Style.Bold = true
			return []RichTextSegment{&textSegment}, nil
		}
	case *ast.ThematicBreak:
		return []RichTextSegment{&SeparatorSegment{}}, nil
	case *ast.Link:
		link, _ := url.Parse(string(t.Destination))
		text := forceIntoText(source, n)
		return []RichTextSegment{&HyperlinkSegment{Text: text, URL: link, Style: TextStyleLink}}, nil
	case *ast.CodeSpan: //代码块
		text := forceIntoText(source, n)
		return []RichTextSegment{&TextSegment{Style: TextStyleCode, Text: text}}, nil
	case *ast.CodeBlock, *ast.FencedCodeBlock: //代码段
		var data []byte
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			data = append(data, line.Value(source)...)
		}
		if len(data) == 0 {
			return nil, nil
		}
		if data[len(data)-1] == '\n' {
			data = data[:len(data)-1]
		}
		return []RichTextSegment{&TextSegment{Style: TextStyleCodeBlock, Text: string(data)}}, nil
	case *ast.Emphasis:
		text := string(forceIntoText(source, n))
		switch t.Level {
		case 2:
			return []RichTextSegment{&TextSegment{Style: TextStyleStrong, Text: text}}, nil
		default:
			return []RichTextSegment{&TextSegment{Style: TextStyleItalic, Text: text}}, nil
		}
	case *ast.Text:
		text := string(t.Value(source))
		if text == "" {
			// These empty text elements indicate single line breaks after non-text elements in goldmark.
			return []RichTextSegment{&TextSegment{Style: TextStyleDefault, Text: " "}}, nil
		}
		if blockquote {
			return []RichTextSegment{&TextSegment{Style: TextStyleNote, Text: text}}, nil
		}
		//需要换行
		_, ok := n.NextSibling().(*ast.Paragraph)
		fmt.Println(ok)
		if t.SoftLineBreak() || ok {
			return []RichTextSegment{
				&TextSegment{Style: TextStyleDefault, Text: text},
				&TextSegment{Style: TextStyleEnter, Text: ""},
			}, nil
		}
		return []RichTextSegment{&TextSegment{Style: TextStyleDefault, Text: text}}, nil
	case *ast.Blockquote:
		return renderChildren(source, n, true)
	case *ast.Image:
		return parseMarkdownImage(t), nil
	}
	return nil, nil
}
func parseMarkdownImage(t *ast.Image) []RichTextSegment {
	return []RichTextSegment{&ImageSegment{
		Path:  string(t.Destination),
		Title: string(t.Title),
		Width: 960,
	},
	}
}

func renderChildren(source []byte, n ast.Node, blockquote bool) ([]RichTextSegment, error) {
	children := make([]RichTextSegment, 0, n.ChildCount())
	for childCount, child := n.ChildCount(), n.FirstChild(); childCount > 0; childCount-- {
		segs, err := renderNode(source, child, blockquote)
		if err != nil {
			return children, err
		}
		children = append(children, segs...)
		child = child.NextSibling()
	}
	return children, nil
}

func forceIntoText(source []byte, n ast.Node) string {
	texts := make([]string, 0)
	ast.Walk(n, func(n2 ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := n2.(type) {
			case *ast.Text:
				texts = append(texts, string(t.Value(source)))
			}
		}
		return ast.WalkContinue, nil
	})
	return strings.Join(texts, " ")
}

func forceIntoHeadingText(source []byte, n ast.Node) string {
	var text strings.Builder
	ast.Walk(n, func(n2 ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := n2.(type) {
			case *ast.Text:
				text.Write(t.Value(source))
			}
		}
		return ast.WalkContinue, nil
	})
	return text.String()
}
