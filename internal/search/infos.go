package search

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/text/width"
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
	// 获取每个格子最长的一格
	var mSearchCode, mName, mAbbr, mExchange, mClassify int
	for _, info := range s.StockInfos {
		updateMaxWidth(calcWidth(info.SearchCode()), &mSearchCode)
		updateMaxWidth(calcWidth(info.Name), &mName)
		updateMaxWidth(calcWidth(info.Abbr), &mAbbr)
		updateMaxWidth(calcWidth(info.Exchange), &mExchange)
		updateMaxWidth(calcWidth(info.Classify), &mClassify)
	}
	const additionWidth = 5
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: mSearchCode + additionWidth},
		{Number: 2, WidthMin: mName + additionWidth},
		{Number: 3, WidthMin: mAbbr + additionWidth},
		{Number: 4, WidthMin: mExchange + additionWidth},
		{Number: 5, WidthMin: mClassify + additionWidth},
	})
	for i, info := range s.StockInfos {
		rows[i] = table.Row{info.SearchCode(), info.Name, info.Abbr, info.Exchange, info.Classify}
	}
	tw.AppendRows(rows)
	tw.SetAutoIndex(true)
	tw.SetCaption("\n\n")
	tw.SetStyle(table.StyleColoredRedWhiteOnBlack)
	return tw.Render()
}

// 计算宽度相关
func calcWidth(s string) (w int) {
	for _, r := range s {
		prop := width.LookupRune(r)
		if prop.Kind() == width.EastAsianFullwidth || prop.Kind() == width.EastAsianWide {
			w += 2
		} else {
			w += 1
		}
	}
	return
}

func updateMaxWidth(curr int, target *int) {
	if curr > *target {
		*target = curr
	}
}
