package stock_data

import (
	"fmt"
	"strings"
)

type Tracker interface {
	tryInitTracker(s *StockData) bool
	evaluate(s *StockData)
	String() string
}

type Trackers []Tracker

func NewTrackers(s *StockData) (ret Trackers) {
	fullList := Trackers{
		&trackerContinuous{},
		&trackerPriceDiff{},
		&trackerPercentDiff{},
		&trackerRealTimePercentDiff{},
		&trackerTargetHighPrice{},
		&trackerTargetLowPrice{},
	}
	for _, tracker := range fullList {
		t := tracker.tryInitTracker(s)
		if t {
			ret = append(ret, tracker)
		}
	}
	if len(ret) == 0 {
		s.Shout("未配置监控器", "")
	} else {
		var msgs []string
		for _, tracker := range ret {
			msgs = append(msgs, tracker.String())
		}
		s.Shout(
			fmt.Sprintf("开始监控%v个指标", len(ret)),
			strings.Join(msgs, "\n"),
		)
	}
	return
}

func (t Trackers) Evaluate(s *StockData) {
	for _, tracker := range t {
		tracker.evaluate(s)
	}
}
