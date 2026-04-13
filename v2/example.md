# Markdown 转图片示例

这是一个基于 **goldmark** 和 `FloatTech/gg` 的纯 Go Markdown 图片渲染器。

它支持以下常见语法：

- 标题
- 段落与自动换行
- **粗体**、*斜体*、`行内代码`
- 有序列表与无序列表
- 引用块
- 代码块
- [链接文本](https://github.com/FloatTech/gg)

> 设计重点不是“把 AST 直接画出来”，而是先做轻量结构转换，再统一布局，这样性能和可维护性都更好。

## 代码块

```go
func RenderMarkdownToPNG(src []byte, output string) error {
	r, err := renderer.New(renderer.Options{
		Theme: renderer.DefaultTheme(1200),
	})
	if err != nil {
		return err
	}
	return r.RenderToFile(src, output)
}
```

1. 先解析 Markdown
2. 再根据宽度完成布局
3. 最后一次性绘制到 PNG
