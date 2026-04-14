package font

import _ "embed"

// Name 表示内置字体资源名。
type Name string

//go:embed MaokenZhuyuanTi.ttf
var TTF []byte
