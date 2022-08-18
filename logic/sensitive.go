package logic

import (
	"strings"
	"ws/global"
)

func FilterSensitive(msg string) string {
	for _, word := range global.SensitiveWords {
		// 把符合条件的内容全部替换成**
		msg = strings.ReplaceAll(msg, word, "**")
	}

	return msg
}
