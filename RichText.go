package marktoimage

import (
	"image/color"

	"github.com/FloatTech/gg"
	"golang.org/x/image/font"
)

// 传参接口
type RichTextSegment interface {
	Inline() bool //绘制本组件是否换行
	SetParent(*RichText)
	Draw()
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
	Block     bool //块,用于存储后面元素是否可以使用Inline
}

// RichText 表示富文本
type RichText struct {
	Segments   []RichTextSegment
	NextInline bool
	Size       float64 //用于保存当前字体大小
	Cov        *gg.Context
	gg.Point   //用于定位画笔
	Config
}
type Config struct {
	FontData        [][]byte //用于全局的字体读取
	TopMargin       float64  //上边距留白
	TextMargin      float64  //左右留白
	Width           float64  //宽度
	Height          float64  //高度
	LineHeight      float64  //行高
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
		r.NextInline = r.Segments[k+1].Inline()
		r.Segments[k].Draw()
		//检查余量
		if r.Height-r.Y < 200 {
			r.Expansion()
		}
	}
	r.NextInline = false
	r.Segments[len(r.Segments)-1].Draw()
}

func (r *RichText) SetFontSize(t TextStyle) {
	if t.FontFace != nil {
		r.Cov.SetFontFace(t.FontFace)
		return
	}
	if t.Size <= 0 || r.Size == t.Size {
		return
	}
	r.Size = t.Size
	r.Cov.ParseFontFace(r.FontData[0], t.Size)
}

func (r *RichText) Cut() {
	if r.Cov.Height() < int(r.Y) {
		return
	}
	r.Y += r.LineHeight
	newDC := gg.NewContext(r.Cov.W(), int(r.Y))
	newDC.DrawImage(r.Cov.Image(), 0, 0)
	r.Cov = newDC
}

// 扩充画布高度
func (r *RichText) Expansion() {
	newDC := gg.NewContext(r.Cov.W(), int(r.Config.Height)*2)
	newDC.DrawImage(r.Cov.Image(), 0, 0)
	r.Cov = newDC
}

var (
	//默认文本
	TextStyleDefault = TextStyle{Inline: true, Block: true}
	//换行文本
	TextStyleBlockquote = TextStyle{Inline: false, Block: true}
	//标准块
	TextStyleParagraph = TextStyle{Inline: false, Size: 40, Block: true}
	//用于标题
	TextStyleHead = TextStyle{Inline: false, Bold: true, Size: 60, Block: true}
	//用于子标题
	TextStyleSubHead = TextStyle{Inline: false, Bold: true, Size: 50, Block: true}
	//用于List的序号样式
	TextStyleList = TextStyle{Inline: false, Bold: true}
	//用于link
	TextStyleLink = TextStyle{Inline: true, Color: color.RGBA{0, 0, 255, 255}}
	//用于code
	TextStyleCode = TextStyle{Inline: true}
	//strong
	TextStyleStrong = TextStyle{Inline: true, Bold: true}
	//Italic
	TextStyleItalic = TextStyle{Inline: true, Italic: true}
)
