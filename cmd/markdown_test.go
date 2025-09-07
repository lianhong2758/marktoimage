package markdown_test

import (
	"fmt"
	"image/color"
	"marktoimage"
	"os"
	"testing"
)

func TestMarkdown(t *testing.T) {
	data, err := os.ReadFile("MaokenZhuyuanTi.ttf")
	if err != nil {
		panic(err)
	}
	rt := marktoimage.NewRichText(
		marktoimage.Config{
			FontData:   [][]byte{data},
			Width:      1080,
			Height:     1920,
			LineHeight: 8,
			TextMargin: 20,
			TopMargin:  10,
		})
	md, _ := os.ReadFile("test.md")
	rt.ParseMarkdown(string(md))
	rt.Draw()
	for _, v := range rt.Segments {
		if l, ok := v.(*marktoimage.TextSegment); ok {
			te := l.Text
			fmt.Println(*l, len(te))
		}
	}
	rt.Cut()
	rt.Cov.SavePNG("md.png")
}

func TestRichText(t *testing.T) {
	data, err := os.ReadFile("regular.ttf")
	if err != nil {
		panic(err)
	}
	rt := marktoimage.NewRichText(
		marktoimage.Config{
			FontData:   [][]byte{data},
			Width:      1080,
			Height:     1920,
			LineHeight: 8,
			TextMargin: 20,
			TopMargin:  10,
		})
	rt.AppendSegment(
		&marktoimage.TextSegment{Style: marktoimage.TextStyleDefault, Text: "abcdefg"},
		&marktoimage.TextSegment{Style: marktoimage.TextStyle{Color: color.RGBA{255, 0, 0, 255}, Inline: false},
			Text: "abcdefg"},
		&marktoimage.SeparatorSegment{},
		&marktoimage.TextSegment{Style: marktoimage.TextStyle{Inline: false}, Text: "abcdefg"},
		&marktoimage.SeparatorSegment{},
		&marktoimage.TextSegment{Style: marktoimage.TextStyle{Inline: false}, Text: "abcdefg"},
		&marktoimage.ListSegment{Ordered: true, Items: []marktoimage.RichTextSegment{
			&marktoimage.TextSegment{Style: marktoimage.TextStyleBlockquote, Text: "abc"},
			&marktoimage.TextSegment{Style: marktoimage.TextStyleDefault, Text: "efg"},
			&marktoimage.TextSegment{Style: marktoimage.TextStyleDefault, Text: "abdefg"},
		}},
		&marktoimage.ImageSegment{
			Path:  "图片1.jpg",
			Title: "Win11",
			Width: 960,
		},
		&marktoimage.HyperlinkSegment{Style: marktoimage.TextStyleLink, Text: "假如这是一个链接"},
		&marktoimage.TextSegment{Style: marktoimage.TextStyleDefault, Text: "这是正常文本"},
	)
	rt.Draw()
	rt.Cut()
	rt.Cov.SavePNG("test.png")
}
