package marktoimage

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark/ast"
	mdtext "github.com/yuin/goldmark/text"
)

type blockKind int

const (
	blockParagraph blockKind = iota
	blockHeading
	blockBlockquote
	blockList
	blockCode
	blockRule
)

type document struct {
	Blocks []block
}

// block 是布局阶段使用的轻量块级结构。
// 它只保留渲染真正需要的信息，避免把完整 AST 传来传去。
type block struct {
	Kind    blockKind
	Level   int
	Info    string
	Inlines []inlineSpan
	Text    string
	Items   []listItem
	Blocks  []block
	Ordered bool
	Start   int
	Tight   bool
}

type listItem struct {
	Blocks []block
}

type inlineSpan struct {
	Text       string
	Bold       bool
	Italic     bool
	Code       bool
	Link       bool
	ForceBreak bool
}

func (r *Renderer) parse(source []byte) (*document, error) {
	root := r.md.Parser().Parse(mdtext.NewReader(source))
	return &document{Blocks: r.parseBlocks(root, source)}, nil
}

func (r *Renderer) parseBlocks(parent ast.Node, source []byte) []block {
	var blocks []block

	for node := parent.FirstChild(); node != nil; node = node.NextSibling() {
		switch n := node.(type) {
		case *ast.Heading:
			blocks = append(blocks, block{
				Kind:    blockHeading,
				Level:   n.Level,
				Inlines: r.parseInlines(n, source, inlineSpan{}),
			})
		case *ast.TextBlock:
			// tight list 中的正文通常会变成 TextBlock，而不是 Paragraph。
			// 这里把它当作普通段落处理，避免列表项正文被漏掉。
			blocks = append(blocks, block{
				Kind:    blockParagraph,
				Inlines: r.parseInlines(n, source, inlineSpan{}),
			})
		case *ast.Paragraph:
			blocks = append(blocks, block{
				Kind:    blockParagraph,
				Inlines: r.parseInlines(n, source, inlineSpan{}),
			})
		case *ast.Blockquote:
			blocks = append(blocks, block{
				Kind:   blockBlockquote,
				Blocks: r.parseBlocks(n, source),
			})
		case *ast.List:
			blocks = append(blocks, r.parseList(n, source))
		case *ast.FencedCodeBlock:
			blocks = append(blocks, block{
				Kind: blockCode,
				Text: readBlockLines(n, source),
				Info: string(bytes.TrimSpace(n.Language(source))),
			})
		case *ast.CodeBlock:
			blocks = append(blocks, block{
				Kind: blockCode,
				Text: readBlockLines(n, source),
			})
		case *ast.ThematicBreak:
			blocks = append(blocks, block{Kind: blockRule})
		default:
			if node.HasChildren() {
				blocks = append(blocks, r.parseBlocks(node, source)...)
			}
		}
	}

	return blocks
}

func (r *Renderer) parseList(n *ast.List, source []byte) block {
	list := block{
		Kind:    blockList,
		Ordered: n.IsOrdered(),
		Start:   n.Start,
		Tight:   n.IsTight,
	}

	for itemNode := n.FirstChild(); itemNode != nil; itemNode = itemNode.NextSibling() {
		item, ok := itemNode.(*ast.ListItem)
		if !ok {
			continue
		}
		list.Items = append(list.Items, listItem{
			Blocks: r.parseBlocks(item, source),
		})
	}

	return list
}

func (r *Renderer) parseInlines(parent ast.Node, source []byte, inherited inlineSpan) []inlineSpan {
	var spans []inlineSpan

	for node := parent.FirstChild(); node != nil; node = node.NextSibling() {
		switch n := node.(type) {
		case *ast.Text:
			text := string(n.Value(source))
			if text != "" {
				span := inherited
				span.Text = text
				spans = append(spans, span)
			}
			if n.SoftLineBreak() {
				span := inherited
				span.Text = " "
				spans = append(spans, span)
			}
			if n.HardLineBreak() {
				span := inherited
				span.ForceBreak = true
				spans = append(spans, span)
			}
		case *ast.String:
			if len(n.Value) > 0 {
				span := inherited
				span.Text = string(n.Value)
				spans = append(spans, span)
			}
		case *ast.Emphasis:
			next := inherited
			if n.Level >= 1 {
				next.Italic = true
			}
			if n.Level >= 2 {
				next.Bold = true
			}
			spans = append(spans, r.parseInlines(n, source, next)...)
		case *ast.CodeSpan:
			next := inherited
			next.Code = true
			spans = append(spans, r.parseInlines(n, source, next)...)
		case *ast.Link:
			next := inherited
			next.Link = true
			spans = append(spans, r.parseInlines(n, source, next)...)
		case *ast.AutoLink:
			next := inherited
			next.Link = true
			next.Text = string(n.Label(source))
			spans = append(spans, next)
		case *ast.Image:
			alt := strings.TrimSpace(r.collectPlainText(n, source))
			if alt == "" {
				alt = "图片"
			}
			span := inherited
			span.Text = "[图片: " + alt + "]"
			spans = append(spans, span)
		default:
			if node.HasChildren() {
				spans = append(spans, r.parseInlines(node, source, inherited)...)
			}
		}
	}

	return spans
}

func (r *Renderer) collectPlainText(parent ast.Node, source []byte) string {
	var sb strings.Builder

	for node := parent.FirstChild(); node != nil; node = node.NextSibling() {
		switch n := node.(type) {
		case *ast.Text:
			sb.Write(n.Value(source))
			if n.SoftLineBreak() || n.HardLineBreak() {
				sb.WriteByte(' ')
			}
		case *ast.String:
			sb.Write(n.Value)
		case *ast.AutoLink:
			sb.Write(n.Label(source))
		default:
			if node.HasChildren() {
				sb.WriteString(r.collectPlainText(node, source))
			}
		}
	}

	return sb.String()
}

func readBlockLines(node interface{ Lines() *mdtext.Segments }, source []byte) string {
	lines := node.Lines()
	var parts []string
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		parts = append(parts, string(segment.Value(source)))
	}
	return strings.Join(parts, "")
}
