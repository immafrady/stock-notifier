package stock_data

import "fmt"

type trackerTargetLowPrice struct {
	flag   int
	target float64
	price  string
}

func (t *trackerTargetLowPrice) tryInitTracker(s *StockData) bool {
	if s.Config.TargetLowPrice > 0 {
		t.target = s.Config.TargetLowPrice
		t.price = s.ApiData.ParsePrice(s.Config.TargetHighPrice)
		return true
	}
	return false
}

func (t *trackerTargetLowPrice) evaluate(s *StockData) {
	if s.ApiData.Current < t.target && t.flag != 2 {
		t.flag = 2
		s.Shout(
			fmt.Sprintf("⬇️跌破目标低位%s :(", t.price),
			"",
		)
	} else if s.ApiData.Current == t.target && t.flag != 1 {
		t.flag = 1
		s.Shout(
			fmt.Sprintf("⏺️️达到目标低位%s", t.price),
			"",
		)
	} else if s.ApiData.Current > t.target && t.flag != 0 {
		t.flag = 0
		s.Shout(
			fmt.Sprintf("⬆️超越目标低位%s :)", t.price),
			"",
		)
	}
}

func (t *trackerTargetLowPrice) String() string {
	return fmt.Sprintf("监控低位价格: %s", t.price)
}
