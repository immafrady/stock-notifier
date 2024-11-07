package search

import (
	"fmt"
	"github.com/immafrady/stock-notifier/utils"
	"strings"
)

// StockInfo 单个股票信息
type StockInfo struct {
	Exchange string // 股票市场
	Code     string // 股票编码
	Name     string // 股票名称
	Abbr     string // 股票简写
	Classify string // 股票分类
}

// SearchCode 查询代码
func (s StockInfo) SearchCode() string {
	if s.Exchange == "us" {
		str := s.Exchange
		str += strings.ToUpper(s.Code)
		return strings.Split(str, ".")[0]
	} else {
		return s.Exchange + s.Code
	}

}

func NewStockInfo(match []string) *StockInfo {
	if len(match) == 6 {
		return &StockInfo{
			Exchange: match[1],
			Code:     match[2],
			Name:     utils.ParseUnicode(match[3]),
			Abbr:     match[4],
			Classify: match[5],
		}
	} else {
		return nil
	}
}

func (s StockInfo) String() string {
	return fmt.Sprintf("查询代码: %s -- 名称: %s -- 缩写: %s -- 分类: %s \n", s.SearchCode(), s.Name, s.Abbr, s.Classify)
}
