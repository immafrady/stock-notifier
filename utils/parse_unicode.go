package utils

import (
	"log"
	"strconv"
)

// ParseUnicode 转化unicode字符
func ParseUnicode(str string) string {
	ret, err := strconv.Unquote(`"` + str + `"`)
	if err != nil {
		log.Println("[error]转化unicode字符失败: ", err.Error())
		return str
	} else {
		return ret
	}
}
