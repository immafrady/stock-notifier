package search

import (
	"fmt"
	"github.com/immafrady/stock-notifier/utils"
)

// StockInfo 单个股票信息
type StockInfo struct {
	Code     string // 股票编码
	Name     string // 股票名称
	Abbr     string // 股票简写
	Classify string // 股票分类
}

func NewStockInfo(match []string) *StockInfo {
	if len(match) == 6 {
		return &StockInfo{
			Code:     match[1] + match[2],
			Name:     utils.ParseUnicode(match[3]),
			Abbr:     match[4],
			Classify: match[5],
		}
	} else {
		return nil
	}
}

func (s StockInfo) String() string {
	return fmt.Sprintf("查询代码: %s -- 名称: %s -- 缩写: %s -- 分类: %s \n", s.Code, s.Name, s.Abbr, s.Classify)
}
