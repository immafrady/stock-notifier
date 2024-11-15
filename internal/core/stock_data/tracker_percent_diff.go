package stock_data

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type trackerPercentDiff struct {
	flag   bool
	target float64
	diff   string
}

func (t *trackerPercentDiff) tryInitTracker(s *StockData) bool {
	t.diff = s.Config.PercentDiff
	percent := strings.TrimSuffix(t.diff, "%")
	value, err := strconv.ParseFloat(percent, 64)
	if err == nil {
		t.target = value / 100
		if t.target != 0 {
			return true
		}
	}
	return false
}

func (t *trackerPercentDiff) evaluate(s *StockData) {
	l := len(s.PriceLogs)
	if l > 2 {
		prev := s.PriceLogs[l-2].Price
		curr := s.PriceLogs[l-1].Price
		diff := (curr - prev) / prev
		if math.Abs(diff) > t.target {
			s.Shout("当前异动已超监控值", fmt.Sprintf("监控时间段内价格变化为： %0.2f%%\n%s➡️%s", diff*100, s.ApiData.ParsePrice(prev), s.ApiData.ParsePrice(curr)))
		}

	}
}

func (t *trackerPercentDiff) String() string {
	return fmt.Sprintf("监控百分比差异: %s", t.diff)
}
