package marktoimage

import (
	"strings"

	"github.com/FloatTech/gg"
)

func TextWrap(dc *gg.Context, s string, w float64) []string {
	w -= 20 //留出容错区间
	s = strings.ReplaceAll(s, "\t", "      ")
	splits := strings.Split(s, "\n")
	t := make([]string, 0, len(splits))
	for _, line := range splits {
		// 如果整行宽度已经小于等于最大宽度，直接添加
		wt, _ := dc.MeasureString(line)
		if wt <= w {
			t = append(t, line)
			continue
		}

		// 字符级断行
		var currentRunes = make([]rune, 0, 30)

		for _, r := range line {
			currentRunes = append(currentRunes, r)
			testLine := string(currentRunes)
			testWt, _ := dc.MeasureString(testLine)

			if testWt > w {
				t = append(t, testLine)
				currentRunes = make([]rune, 0, 30)
			}
		}
		t = append(t, string(currentRunes))
	}
	return t
}
