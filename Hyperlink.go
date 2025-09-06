package marktoimage

import (
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
	t := TextSegment{Text: h.Text, Style: h.Style}
	t.SetParent(h.parent)
	t.Draw()
}

func (h *HyperlinkSegment) SetParent(r *RichText) {
	h.parent = r
}
