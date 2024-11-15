package stock_data

import (
	"fmt"
	"math"
)

type trackerPriceDiff struct {
	flag   bool
	target float64
}

func (t *trackerPriceDiff) tryInitTracker(s *StockData) bool {
	if s.Config.PriceDiff > 0 {
		t.target = s.Config.PriceDiff
		return true
	}
	return false
}

func (t *trackerPriceDiff) evaluate(s *StockData) {
	l := len(s.PriceLogs)
	if l > 2 {
		prev := s.PriceLogs[l-2].Price
		curr := s.PriceLogs[l-1].Price
		diff := curr - prev
		if math.Abs(diff) > t.target {
			if !t.flag {
				t.flag = true
				s.Shout("当前差价已超监控值", fmt.Sprintf("监控时间段内价格变化为： %s", s.ApiData.ParsePrice(diff)))
			}
		} else {
			t.flag = false
		}
	}
}

func (t *trackerPriceDiff) String() string {
	return fmt.Sprintf("监控价格差异: %v", t.target)

}
