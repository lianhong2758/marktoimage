package marktoimage

import (
	"image/color"
	"net/url"
)

type HyperlinkSegment struct {
	Style TextStyle
	Text  string
	URL   *url.URL

	parent *RichText
}

func (h *HyperlinkSegment) Inline() bool {
	return h.Style.Inline
}
func (h *HyperlinkSegment) Draw() {
	h.Style.Color = h.color()
	t := TextSegment{Text: h.Text, Style: h.Style}
	t.SetParent(h.parent)
	t.Draw()
}

func (h *HyperlinkSegment) SetParent(r *RichText) {
	h.parent = r
}
func (h *HyperlinkSegment) color() color.Color {
	if h.Style.Color == nil {
		return h.parent.Config.Colors.LinkColor
	}
	return h.Style.Color
}
