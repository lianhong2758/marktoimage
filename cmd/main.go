package main

import (
	"fmt"
	"marktoimage"
	"os"
)

func main() {
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
	md, _ := os.ReadFile("test.md")
	rt.ParseMarkdown(string(md))
	rt.Draw()
	for _, v := range rt.Segments {
		if l, ok := v.(*marktoimage.TextSegment); ok {
			fmt.Println(*l)
		}
	}
	rt.Cov.SavePNG("md.png")
}

// func main() {
// 	data, err := os.ReadFile("regular.ttf")
// 	if err != nil {
// 		panic(err)
// 	}
// 	rt := marktoimage.NewRichText(
// 		marktoimage.Config{
// 			FontData:   [][]byte{data},
// 			Width:      1080,
// 			Height:     1920,
// 			LineHeight: 8,
// 			TextMargin: 20,
// 			TopMargin:  10,
// 		})
// 	rt.AppendSegment(
// 		&marktoimage.TextSegment{Style: marktoimage.TextStyleDefault, Text: "abcdefg"},
// 		&marktoimage.TextSegment{Style: marktoimage.TextStyle{Color: color.RGBA{255, 0, 0, 255}, Inline: true},
// 			Text: "abcdefg"},
// 		&marktoimage.SeparatorSegment{},
// 		&marktoimage.TextSegment{Style: marktoimage.TextStyle{Inline: false}, Text: "abcdefg"},
// 		&marktoimage.SeparatorSegment{},
// 		&marktoimage.TextSegment{Style: marktoimage.TextStyle{Inline: false}, Text: "abcdefg"},
// 		&marktoimage.ListSegment{Ordered: true, Items: []marktoimage.RichTextSegment{
// 			&marktoimage.TextSegment{Style: marktoimage.TextStyleDefault, Text: "abc"},
// 			&marktoimage.TextSegment{Style: marktoimage.TextStyleDefault, Text: "efg"},
// 			&marktoimage.TextSegment{Style: marktoimage.TextStyleDefault, Text: "abdefg"},
// 		}},
// 		&marktoimage.ImageSegment{
// 			Path:  "图片1.jpg",
// 			Title: "Win11",
// 			Width: 960,
// 		},
// 		&marktoimage.HyperlinkSegment{Style: marktoimage.TextStyle{Inline: true, Color: color.RGBA{0, 0, 255, 255}}, Text: "假如这是一个链接"},
// 	)
// 	rt.SetFontSize(marktoimage.TextStyle{Inline: true, Size: 40})
// 	rt.Draw()
// 	rt.Cov.SavePNG("test.png")
// }
