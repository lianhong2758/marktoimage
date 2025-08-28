package marktoimage

import (
	"image/color"

	"github.com/FloatTech/gg"
	"golang.org/x/image/font"
)

// 传参接口
type RichTextSegment interface {
	Inline() bool
	SetParent(*RichText)
	Draw(bool)
}

// TextStyle 表示文本样式
type TextStyle struct {
	Inline    bool //是否在一行绘制/是否不要新建一行
	Size      float64
	FontFace  font.Face
	Color     color.Color
	Bold      bool //加粗
	Italic    bool //斜体
	Underline bool //下划线
}

// RichText 表示富文本
type RichText struct {
	Segments []RichTextSegment
	Cov      *gg.Context
	gg.Point //用于定位画笔
	Config
}
type Config struct {
	TopMargin       float64 //上边距留白
	TextMargin      float64 //左右留白
	Width           float64 //宽度
	Height          float64 //高度
	LineHeight      float64 //行高
	DefaultColor    color.Color
	BackgroundColor color.Color
}

func NewRichText(cfg Config) *RichText {
	if cfg.DefaultColor == nil {
		cfg.DefaultColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	}
	if cfg.BackgroundColor == nil {
		cfg.BackgroundColor = color.NRGBA{255, 255, 255, 255}
	}
	cov := gg.NewContext(int(cfg.Width), int(cfg.Height))
	cov.SetRGB(1, 1, 1)
	cov.Clear()
	return &RichText{
		Cov:      cov,
		Segments: []RichTextSegment{},
		Point:    gg.Point{X: 0, Y: cfg.TopMargin},
		Config:   cfg,
	}
}
func (r *RichText) AppendSegment(rs ...RichTextSegment) {
	for i := range rs {
		rs[i].SetParent(r)
	}
	r.Segments = append(r.Segments, rs...)
}

func (r *RichText) Draw() {
	for k := 0; k < len(r.Segments)-1; k++ {
		r.Segments[k].Draw(r.Segments[k+1].Inline())
	}
	r.Segments[len(r.Segments)-1].Draw(false)
}

var (
	//用于List的序号样式
	ListTextStyle = TextStyle{Inline: true, Bold: true}
)
