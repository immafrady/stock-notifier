package config

import (
	"strconv"
	"strings"
)

// TimeString 短的时间戳
type TimeString string

func (t TimeString) Get() (h, m int) {
	strs := strings.Split(string(t), ":")
	if len(strs) != 2 {
		return 0, 0
	} else {
		h, _ = strconv.Atoi(strs[0])
		m, _ = strconv.Atoi(strs[1])
		return
	}
}
