package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lianhong2758/marktoimage/renderer"
)

const sampleMarkdown = `# Markdown 转图片示例

这是一个基于 **goldmark** 和 ` + "`FloatTech/gg`" + ` 的纯 Go Markdown 图片渲染器。

它支持以下常见语法：

- 标题
- 段落与自动换行
- **粗体**、*斜体*、` + "`行内代码`" + `
- 有序列表与无序列表
- 引用块
- 代码块
- [链接文本](https://github.com/FloatTech/gg)

> 设计重点不是“把 AST 直接画出来”，而是先做轻量结构转换，再统一布局，这样性能和可维护性都更好。

## 代码块

` + "```go" + `
func RenderMarkdownToPNG(src []byte, output string) error {
	r, err := renderer.New(renderer.Options{
		ThemeName: renderer.ThemeDefault,
		Width:     1200,
	})
	if err != nil {
		return err
	}
	return r.RenderToFile(src, output)
}
` + "```" + `

1. 先解析 Markdown
2. 再根据宽度完成布局
3. 最后一次性绘制到 PNG`

func main() {
	var (
		in    = flag.String("in", "", "Markdown 文件路径，留空则使用内置示例")
		out   = flag.String("out", "output/markdown.png", "输出 PNG 路径")
		width = flag.Int("width", 1200, "图片宽度")
		font  = flag.String("font", "", "TTF 字体文件路径，留空时自动探测常见中文字体位置")
		theme = flag.String("theme", string(renderer.ThemeDefault), "内置主题，可选: "+strings.Join(renderer.ThemeNames(), ", "))
	)
	flag.Parse()

	content := []byte(sampleMarkdown)
	if *in != "" {
		data, err := os.ReadFile(*in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "读取 Markdown 文件失败: %v\n", err)
			os.Exit(1)
		}
		content = data
	}

	themeName, err := renderer.ParseThemeName(*theme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "主题配置无效: %v\n", err)
		os.Exit(1)
	}

	fontPath, autoDetected, err := resolveFontPath(*font)
	if err != nil {
		fmt.Fprintf(os.Stderr, "字体配置无效: %v\n", err)
		os.Exit(1)
	}
	if autoDetected != "" {
		fmt.Fprintf(os.Stderr, "使用自动探测到的字体: %s\n", autoDetected)
	} else if fontPath == "" {
		fmt.Fprintln(os.Stderr, "未找到可用的自定义字体，回退到内置 Go 字体；中文可能无法正常显示，可通过 -font 指定 TTF。")
	}

	r, err := renderer.New(renderer.Options{
		ThemeName: themeName,
		Width:     *width,
		FontPath:  fontPath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化渲染器失败: %v\n", err)
		os.Exit(1)
	}

	output := filepath.Clean(*out)
	if err := r.RenderToFile(content, output); err != nil {
		fmt.Fprintf(os.Stderr, "生成图片失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已生成图片: %s\n", output)
}

func resolveFontPath(raw string) (resolved string, autoDetected string, err error) {
	if strings.TrimSpace(raw) != "" {
		path := filepath.Clean(raw)
		if !fileExists(path) {
			return "", "", fmt.Errorf("font file %q does not exist", path)
		}
		return path, "", nil
	}

	for _, candidate := range autoFontCandidates() {
		if fileExists(candidate) {
			return candidate, candidate, nil
		}
	}

	return "", "", nil
}

func autoFontCandidates() []string {
	return []string{
		"MaokenZhuyuanTi.ttf",
		filepath.Join("..", "MaokenZhuyuanTi.ttf"),
		filepath.Join("..", "cmd", "MaokenZhuyuanTi.ttf"),
		filepath.Join("v2", "MaokenZhuyuanTi.ttf"),
		filepath.Join("cmd", "MaokenZhuyuanTi.ttf"),
		filepath.Join("..", "..", "MaokenZhuyuanTi.ttf"),
		filepath.Join("..", "..", "cmd", "MaokenZhuyuanTi.ttf"),
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
