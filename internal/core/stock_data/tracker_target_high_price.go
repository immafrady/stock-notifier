package stock_data

import "fmt"

type trackerTargetHighPrice struct {
	flag   int
	target float64
	price  string
}

func (t *trackerTargetHighPrice) tryInitTracker(s *StockData) bool {
	if s.Config.TargetHighPrice > 0 {
		t.target = s.Config.TargetHighPrice
		t.price = s.ApiData.ParsePrice(s.Config.TargetHighPrice)
		return true
	}
	return false
}

func (t *trackerTargetHighPrice) evaluate(s *StockData) {
	if s.ApiData.Current > t.target && t.flag != 2 {
		t.flag = 2
		s.Shout(
			fmt.Sprintf("⬆️超越目标高位%s :)", t.price),
			"",
		)
	} else if s.ApiData.Current == t.target && t.flag != 1 {
		t.flag = 1
		s.Shout(
			fmt.Sprintf("⏺️达到目标高位%s", t.price),
			"",
		)
	} else if s.ApiData.Current < t.target && t.flag != 0 {
		t.flag = 0
		s.Shout(
			fmt.Sprintf("⬇️跌破目标高位%s :(", t.price),
			"",
		)
	}
}

func (t *trackerTargetHighPrice) String() string {
	return fmt.Sprintf("监控高位价格: %v", t.price)
}
