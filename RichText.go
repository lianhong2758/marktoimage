package marktoimage

import (
	"image/color"

	"github.com/FloatTech/gg"
	"golang.org/x/image/font"
)

// 传参接口
type RichTextSegment interface {
	Inline() bool //本组件没有结束,继续在本行绘制
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
	Block     bool //是否是块代码,作为一个整体绘制,解析\n
	Base      bool //背景色
}

// RichText 表示富文本
type RichText struct {
	Segments []RichTextSegment
	Size     float64 //用于保存当前字体大小
	Cov      *gg.Context
	gg.Point //用于定位画笔
	Config
}
type Config struct {
	FontData   [][]byte //用于全局的字体读取
	TopMargin  float64  //上边距留白
	TextMargin float64  //左右留白
	Width      float64  //宽度
	Height     float64  //高度
	LineHeight float64  //行高
	Colors     Colors
}
type Colors struct {
	DefaultColor     color.Color //默认字体颜色
	BackgroundColor  color.Color
	CodeBlockBgColor color.Color
	CodeColor        color.Color
	NoteColor        color.Color //引用
	LinkColor        color.Color
}

func NewRichText(cfg Config) *RichText {
	cfg.setDefaultColor()
	cov := gg.NewContext(int(cfg.Width), int(cfg.Height))
	cov.SetColor(cfg.Colors.BackgroundColor)
	cov.Clear()
	return &RichText{
		Cov:      cov,
		Segments: []RichTextSegment{},
		Point:    gg.Point{X: 0, Y: cfg.TopMargin},
		Config:   cfg,
	}
}

func (cfg *Config) setDefaultColor() {
	if cfg.Colors.DefaultColor == nil {
		r, g, b, a := gg.ParseHexColor("#d1d7e0")
		cfg.Colors.DefaultColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	if cfg.Colors.BackgroundColor == nil {
		r, g, b, a := gg.ParseHexColor("#212830")
		cfg.Colors.BackgroundColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	if cfg.Colors.CodeBlockBgColor == nil {
		r, g, b, a := gg.ParseHexColor("#323943")
		cfg.Colors.CodeBlockBgColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	if cfg.Colors.CodeColor == nil {
		r, g, b, a := gg.ParseHexColor("#d1d7e0")
		cfg.Colors.CodeColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	if cfg.Colors.NoteColor == nil {
		r, g, b, a := gg.ParseHexColor("#7a818a")
		cfg.Colors.NoteColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	TextStyleNote.Color = cfg.Colors.NoteColor
	if cfg.Colors.LinkColor == nil {
		r, g, b, a := gg.ParseHexColor("#478be6")
		cfg.Colors.LinkColor = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}

}

func (r *RichText) AppendSegment(rs ...RichTextSegment) {
	for i := range rs {
		rs[i].SetParent(r)
	}
	r.Segments = append(r.Segments, rs...)
}

func (r *RichText) Draw() {
	r.Align()
	for k := 0; k < len(r.Segments)-1; k++ {
		r.Segments[k].Draw()
		//检查余量
		if r.Height-r.Y < 200 {
			r.Expansion()
		}
	}
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
	newDC := gg.NewContext(r.Cov.W(), int(r.Y)+10)
	newDC.DrawImage(r.Cov.Image(), 0, 0)
	r.Cov = newDC
}

// 扩充画布高度
func (r *RichText) Expansion() {
	newDC := gg.NewContext(r.Cov.W(), int(r.Config.Height)*2)
	newDC.DrawImage(r.Cov.Image(), 0, 0)
	r.Cov = newDC
}

// 防止最后一行时.Inline()==true导致Y没有移动到应该切割的位置,需要补充换行
func (r *RichText) Align() {
	if r.Segments[len(r.Segments)-1].Inline() {
		r.AppendSegment(&TextSegment{Style: TextStyleEnter, Text: ""})
	}
}

// 清空richtext,准备解析下一个内容
func (r *RichText) Clear() {
	cov := gg.NewContext(int(r.Width), int(r.Height))
	cov.SetColor(r.Colors.BackgroundColor)
	cov.Clear()
	r.Cov = cov
	r.Segments = []RichTextSegment{}
	r.Point = gg.Point{X: 0, Y: r.TopMargin}
}

var (
	//默认文本
	TextStyleDefault = TextStyle{Inline: true, Size: 40}
	//换行文本
	TextStyleEnter = TextStyle{Inline: false, Size: 40}
	//用于组件间距
	TextStyleParagraph = TextStyle{Inline: false, Size: 20}
	//四级标题
	TextStyleHead = TextStyle{Inline: false, Bold: true, Size: 60}
	//三级标题
	TextStyleSubHead = TextStyle{Inline: false, Bold: true, Size: 50}
	//用于List的序号样式
	TextStyleList = TextStyle{Inline: true, Size: 40}
	//用于link
	TextStyleLink = TextStyle{Inline: true}
	//用于code代码块`abc`
	TextStyleCode = TextStyle{Inline: true, Size: 40, Base: true}
	//用于code代码段```abc\ndef```
	TextStyleCodeBlock = TextStyle{Inline: false, Size: 35, Block: true, Base: true}
	//引用
	//TextStyleNote.Color = cfg.Colors.NoteColor
	TextStyleNote = TextStyle{Inline: false, Size: 35}
	//strong
	TextStyleStrong = TextStyle{Inline: true, Bold: true}
	//Italic
	TextStyleItalic = TextStyle{Inline: true, Italic: true}
)
