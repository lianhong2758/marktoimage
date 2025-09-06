package marktoimage

import (
	"image/color"
)

// TextSegment 表示文本片段,实现RichTextSegment
type TextSegment struct {
	Style TextStyle //文本样式
	Text  string

	parent *RichText
}

// 内联
func (t *TextSegment) Inline() bool {
	return t.Style.Inline
}

// 绘制
func (t *TextSegment) Draw() {
	t.parent.SetFontSize(t.Style)
	t.parent.Cov.SetColor(t.color())
	if t.parent.X == 0 { //如果是开头就加间隔
		t.parent.X += t.parent.TextMargin
	}
	t.parent.Cov.DrawString(t.Text, t.parent.X, t.parent.Y+t.parent.Size)
	ax, _ := t.parent.Cov.MeasureString(t.Text)
	t.parent.X += ax + float64(int(t.parent.LineHeight)/2)
	//nextline
	if !t.Inline() {
		t.parent.X = 0
		t.parent.Y += t.parent.Size + t.parent.LineHeight
	}
}

 

func (t *TextSegment) color() color.Color {
	if t.Style.Color == nil {
		return t.parent.Config.Colors.DefaultColor
	}
	return t.Style.Color
}

func (t *TextSegment) SetParent(r *RichText) {
	t.parent = r
}
