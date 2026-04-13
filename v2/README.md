# mdrender

一个纯 Go 的 Markdown 转图片示例项目：

- Markdown 解析基于 `github.com/yuin/goldmark`
- 绘图基于 `github.com/FloatTech/gg`
- 内置 `default` 和 `github-dark` 两套主题，可通过配置切换

项目目标：

- 尽量遵循 CommonMark 常见语法
- 通过“两阶段布局 + 一次性绘制”提升性能与稳定性
- 保持结构简单，便于继续扩展

快速运行：

```powershell
go run ./cmd/mdrender
```

渲染指定 Markdown 文件：

```powershell
go run ./cmd/mdrender -in ./example.md -out ./output/example.png
```

切换到 GitHub 风格深色主题：

```powershell
go run ./cmd/mdrender -in ./example.md -out ./output/example-dark.png -theme github-dark
```

如果你在代码里初始化渲染器，也可以直接传入主题名：

```go
r, err := renderer.New(renderer.Options{
	ThemeName: renderer.ThemeGitHubDark,
	Width:     1200,
	FontPath:  "MaokenZhuyuanTi.ttf",
})
```
