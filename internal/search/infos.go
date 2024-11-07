package search

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"regexp"
)

// StockInfos 股票信息集群
type StockInfos struct {
	Condition  string
	StockInfos []StockInfo
}

func NewStockInfos(condition string, result string) *StockInfos {
	regex := regexp.MustCompile(`(\w+)~([\w.]+)~([^~]+)~(\w+)~(\w+)`)
	matches := regex.FindAllStringSubmatch(result, -1)

	if len(matches) > 0 {
		infos := &StockInfos{
			Condition: condition,
		}

		for _, match := range matches {
			info := NewStockInfo(match)
			if info != nil {
				infos.StockInfos = append(infos.StockInfos, *info)
			}
		}
		return infos
	} else {
		return nil
	}
}

func (s StockInfos) String() string {
	tw := table.NewWriter()
	tw.SetTitle(s.Condition)
	tw.AppendHeader(table.Row{"查询代码", "名称", "缩写", "所属市场", "分类"})
	rows := make([]table.Row, len(s.StockInfos))
	for i, info := range s.StockInfos {
		rows[i] = table.Row{info.SearchCode(), info.Name, info.Abbr, info.Exchange, info.Classify}
	}
	tw.AppendRows(rows)
	tw.SetAutoIndex(true)
	tw.SetCaption("\n\n")
	tw.SetStyle(table.StyleBold)
	return tw.Render()
}
