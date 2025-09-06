package marktoimage

type SeparatorSegment struct {
	parent *RichText
}

func (SeparatorSegment) Inline() bool {
	return false
}
func (s *SeparatorSegment) Draw() {
	s.parent.Y += s.parent.LineHeight * 3
	s.parent.Cov.SetColor(s.parent.DefaultColor)
	s.parent.Cov.SetLineWidth(3)
	s.parent.Cov.MoveTo(s.parent.X, s.parent.Y)
	s.parent.Cov.LineTo(s.parent.X+s.parent.Width, s.parent.Y)
	s.parent.Cov.Stroke()
	s.parent.X = 0
	s.parent.Y += s.parent.LineHeight*2
}
func (s *SeparatorSegment) SetParent(r *RichText) {
	s.parent = r
}
