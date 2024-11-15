package stock_data

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type trackerRealTimePercentDiff struct {
	flag   bool
	target float64
	diff   string
}

func (t *trackerRealTimePercentDiff) tryInitTracker(s *StockData) bool {
	t.diff = s.Config.RealTimePercentDiff
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

func (t *trackerRealTimePercentDiff) evaluate(s *StockData) {
	rt := s.ApiData.RealTime
	c := s.ApiData.Current
	diff := (c - rt) / rt
	if math.Abs(diff) > t.target {
		var status string
		if diff > 0 {
			status = "溢价"
		} else {
			status = "折价"
		}
		s.Shout(
			fmt.Sprintf("发生%s", status),
			fmt.Sprintf("当前实际净值为%v, %s率为: %0.2f%%", s.ApiData.ParsePrice(s.ApiData.RealTime), status, math.Abs(diff*100)),
		)
	}
}

func (t *trackerRealTimePercentDiff) String() string {
	return fmt.Sprintf("监控溢价/折价: %s", t.diff)
}
