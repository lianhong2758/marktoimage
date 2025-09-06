package marktoimage

import "strconv"

type ListSegment struct {
	Items   []RichTextSegment
	Ordered bool

	segments []RichTextSegment

	parent *RichText
}

func (ListSegment) Inline() bool {
	return false
}
func (l *ListSegment) Draw() {
	l.segments = l.Segments()
	for i, v := range l.segments {
		if l.parent.X == 0 { //如果是开头就加间隔
			l.parent.X += l.parent.TextMargin * 2
		}
		v.Draw()
		//奇数块内换行,如果已经换行则忽略
		if i%2 == 1 && v.Inline() {
			l.parent.X = l.parent.TextMargin * 2
			l.parent.Y += l.parent.Size + l.parent.LineHeight
		}
	}
	//结尾换行
	if l.segments[len(l.segments)-1].Inline() {
		l.parent.X = 0
		//l.parent.Y += l.parent.Size + l.parent.LineHeight
	}
}

func (l *ListSegment) SetParent(r *RichText) {
	l.parent = r
}

//  更新需要绘制的文本
func (l *ListSegment) Segments() []RichTextSegment {
	out := make([]RichTextSegment, len(l.Items)*2)
	for i, in := range l.Items {
		txt := "• "
		if l.Ordered {
			txt = strconv.Itoa(i+1) + "."
		}
		out[i*2] = &TextSegment{Text: txt + " ", Style: TextStyleList}
		out[i*2].SetParent(l.parent)
		out[i*2+1] = in
		out[i*2+1].SetParent(l.parent)
	}
	return out
}
