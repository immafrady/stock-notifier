package stock_data

import (
	"fmt"
	"strings"
)

type trackerContinuous struct {
	target int
}

func (t *trackerContinuous) tryInitTracker(s *StockData) bool {
	if s.Config.Continuous > 2 {
		t.target = s.Config.Continuous
		return true
	}
	return false
}

func (t *trackerContinuous) evaluate(s *StockData) {
	l := len(s.PriceLogs)
	if l == t.target {

		isIncreasing := true
		isDecreasing := true

		for i := 1; i < len(s.PriceLogs); i++ {
			if s.PriceLogs[i].Price > s.PriceLogs[i-1].Price {
				isDecreasing = false // 如果出现上涨，则不是下跌
			} else if s.PriceLogs[i].Price < s.PriceLogs[i-1].Price {
				isIncreasing = false // 如果出现下跌，则不是上涨
			} else {
				return // 平价直接再见
			}
		}

		var trend string
		if isIncreasing {
			trend = "上涨"
		} else if isDecreasing {
			trend = "下跌"
		} else {
			return
		}
		first := s.PriceLogs[0].Price
		last := s.PriceLogs[len(s.PriceLogs)-1].Price
		diff := last - first
		rate := diff / first
		list := make([]string, len(s.PriceLogs))
		for i, p := range s.PriceLogs {
			list[i] = s.ApiData.ParsePrice(p.Price)
		}
		s.Shout(
			fmt.Sprintf("连续%v次记录呈现单边%s走势", t.target, trend),
			fmt.Sprintf("监控时间段内价格变化为： %s(%0.2f%%)\n%v", s.ApiData.ParsePrice(diff), rate, strings.Join(list, "➡️")),
		)
	}
}

func (t *trackerContinuous) String() string {
	return fmt.Sprintf("监控价格连续单向变化: %v次", t.target)
}
