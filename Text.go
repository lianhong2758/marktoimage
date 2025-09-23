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
	switch t.Style.Block {
	case false:
		w, _ := t.parent.Cov.MeasureString(t.Text)
		if w > float64(t.parent.Cov.W())-2*t.parent.Config.TextMargin { //溢出,需要分行
			lines := TextWrap(t.parent.Cov, t.Text, float64(t.parent.Cov.W())-2*t.parent.Config.TextMargin)
			for i, line := range lines {
				t1 := TextSegment{Text: line, Style: t.Style}
				t1.SetParent(t.parent)
				t1.Draw()
				if i < len(lines)-1 && t.Inline() {
					t1 = TextSegment{Style: TextStyleParagraph, Text: ""}
					t1.SetParent(t.parent)
					t1.Draw()
				}
			}
		}
		//绘制背景
		if t.Style.Base {
			//绘制一个
			w, _ = t.parent.Cov.MeasureString(t.Text)
			t.parent.Cov.SetColor(t.parent.Colors.CodeBlockBgColor)
			t.parent.Cov.DrawRoundedRectangle(t.parent.X, t.parent.Y, w, t.parent.Cov.FontHeight()*1.8, 5)
			t.parent.Cov.Fill()
		}
		t.parent.Cov.SetColor(t.color())
		t.parent.Cov.DrawString(t.Text, t.parent.X, t.parent.Y+t.parent.Size)
		ax, _ := t.parent.Cov.MeasureString(t.Text)
		t.parent.X += ax + float64(int(t.parent.LineHeight)/2)
		if !t.Inline() {
			t.parent.X = 0
			t.parent.Y += t.parent.Size + t.parent.LineHeight
		}
	case true:
		//绘制本文
		//预留两个空格的前缀距离
		w, _ := t.parent.Cov.MeasureString("  ")
		lines := TextWrap(t.parent.Cov, t.Text, float64(t.parent.Cov.W())-2*t.parent.Config.TextMargin-w)
		for I := range lines {
			lines[I] = "  " + lines[I]
		}
		h := float64(len(lines)) * t.parent.Cov.FontHeight() * 1.5
		h -= 0.5 * t.parent.Cov.FontHeight()
		y := t.parent.Y + t.parent.Size
		if t.Style.Base {
			t.parent.Cov.SetColor(t.parent.Colors.CodeBlockBgColor)
			t.parent.Cov.DrawRoundedRectangle(t.parent.X, t.parent.Y, float64(t.parent.Cov.W())-2*t.parent.Config.TextMargin, h+t.parent.Cov.FontHeight(), 5)
			t.parent.Cov.Fill()
		}
		t.parent.Cov.SetColor(t.color())
		for _, line := range lines {
			t.parent.Cov.DrawString(line, t.parent.X, y)
			y += t.parent.Cov.FontHeight() * 1.5
		}
		t.parent.X = 0
		t.parent.Y += h + t.parent.Cov.FontHeight()
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
