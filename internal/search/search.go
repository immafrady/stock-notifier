package search

import (
	"fmt"
	"github.com/immafrady/stock-notifier/utils"
	"log"
)

func Find(condition string) {
	str := utils.Request("https://smartbox.gtimg.cn/s3/?v=2&t=all&c=1&q=" + condition)

	infos := NewStockInfos(condition, str)

	if infos != nil {
		fmt.Println(infos)
	} else {
		log.Println("【" + condition + "】 此条件无搜索结果")
	}
}
