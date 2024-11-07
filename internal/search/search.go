package search

import (
	"fmt"
	"github.com/immafrady/stock-notifier/utils"
	"log"
)

func find(condition string) {
	str := utils.Request("https://smartbox.gtimg.cn/s3/?v=2&t=all&c=1&q=" + condition)

	infos := NewStockInfos(condition, str)

	if infos != nil {
		fmt.Println(infos)
	} else {
		log.Println("【" + condition + "】 此条件无搜索结果")
	}
}

func FindAll(conditions []string) {
	for _, condition := range conditions {
		find(condition)
	}
	fmt.Println("将【查询代码】复制到配置文件中，开启自动盯盘")
}
