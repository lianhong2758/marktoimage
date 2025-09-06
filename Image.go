package marktoimage

import (
	"bufio"
	"bytes"
	"image"
	"net/http"
	"strings"

	"github.com/FloatTech/gg"
	"github.com/nfnt/resize"
)

type ImageSegment struct {
	Path   string //支持的有相对路径,绝对路径,网络路径,base64
	Title  string
	Data   []byte      //也可以不写path构建
	Image  image.Image //也可以直接传入图片
	Width  int
	Height int

	parent *RichText
}

func (i *ImageSegment) Inline() bool {
	return false
}

func (i *ImageSegment) SetParent(r *RichText) {
	i.parent = r
}

func (i *ImageSegment) Draw() {
	i.InitImage()
	i.resize()
	if i.Image == nil {
		t := TextSegment{Text: i.Title, Style: TextStyle{Inline: false, Color: i.parent.Colors.DefaultColor}}
		t.SetParent(i.parent)
		t.Draw()
	}
	i.parent.Cov.DrawImageAnchored(i.Image, i.parent.Cov.W()/2, int(i.parent.Y+i.parent.LineHeight*2), 0.5, 0)
	i.parent.X = i.parent.Config.TextMargin
	i.parent.Y += i.parent.LineHeight*3 + float64(i.Image.Bounds().Dy())
}

func (i *ImageSegment) InitImage() {
	if i.Image != nil {
		return
	}
	if len(i.Data) > 0 {
		var err error
		i.Image, _, err = image.Decode(bytes.NewReader(i.Data))
		if err != nil {
			i.Title += " [image] "
		}
		return
	}
	if i.Path == "" {
		i.Title += " [image] "
		return
	}
	if strings.HasPrefix(i.Path, "http://") || strings.HasPrefix(i.Path, "https://") {
		resp, err := http.Get(i.Path)
		if err != nil {
			i.Title += " [image] "
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			i.Title += " [image] "
			return
		}
		i.Image, _, err = image.Decode(bufio.NewReader(resp.Body))
		if err != nil {
			i.Title += " [image] "
			return
		}
	} else {
		//本地图片
		var err error
		i.Image, err = gg.LoadImage(i.Path)
		if err != nil {
			i.Title += " [image] "
			return
		}
	}
}

func (i *ImageSegment) resize() {
	if (i.Width == 0 && i.Height == 0) || i.Image == nil || i.Image.Bounds().Dx() <= 1000 {
		return
	}
	i.Image = resize.Resize(uint(i.Width), uint(i.Height), i.Image, resize.Bicubic)
}
