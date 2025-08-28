package main

import (
	"fmt"
	"image/color"
	"marktoimage"

	"github.com/FloatTech/gg"
)

func main() {
	face, err := gg.LoadFontFace("regular.ttf", 40)
	if err != nil {
		fmt.Println(err)
		return
	}
	sy := marktoimage.TextStyle{Inline: true, FontFace: face}
	rt := marktoimage.NewRichText(
		marktoimage.Config{
			Width:      500,
			Height:     500,
			LineHeight: 5,
			TextMargin: 20,
			TopMargin:  10,
		})
	rt.AppendSegment(
		&marktoimage.TextSegment{Style: sy, Text: "abcdefg"},
		&marktoimage.TextSegment{Style: marktoimage.TextStyle{FontFace: face, Color: color.RGBA{255, 0, 0, 255}, Inline: true},
			Text: "abcdefg"},
		&marktoimage.SeparatorSegment{},
		&marktoimage.TextSegment{Style: marktoimage.TextStyle{FontFace: face, Inline: false}, Text: "abcdefg"},
		&marktoimage.SeparatorSegment{},
		&marktoimage.TextSegment{Style: marktoimage.TextStyle{FontFace: face, Inline: false}, Text: "abcdefg"},
		&marktoimage.ListSegment{Ordered: true,Items: []marktoimage.RichTextSegment{
			&marktoimage.TextSegment{Style: sy, Text: "abc"},
			&marktoimage.TextSegment{Style: sy, Text: "efg"},
			&marktoimage.TextSegment{Style: sy, Text: "abdefg"},
		}},
	)
	rt.Draw()
	rt.Cov.SavePNG("test.png")
}
